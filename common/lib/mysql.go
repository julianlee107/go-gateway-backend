package lib

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/e421083458/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

func InitDBPool(path string) error {
	DBConfMap := &MysqlMapConf{}
	err := ParseConfig(path, DBConfMap)
	if err != nil {
		return err
	}
	if len(DBConfMap.Map) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}
	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}
	for confName, dbConf := range DBConfMap.Map {
		dbpool, err := sql.Open("mysql", dbConf.DataSourceName)
		if err != nil {
			return err
		}
		dbpool.SetMaxIdleConns(dbConf.MaxIdleConn)
		dbpool.SetMaxOpenConns(dbConf.MaxOpenConn)
		dbpool.SetConnMaxLifetime(time.Duration(dbConf.MaxConnLifeTime) * time.Second)
		// test mysql connection
		err = dbpool.Ping()
		if err != nil {
			return err
		}

		dbGorm, err := gorm.Open("mysql", dbConf.DataSourceName)
		if err != nil {
			return err
		}

		dbGorm.SingularTable(true)
		err = dbGorm.DB().Ping()
		if err != nil {
			return err
		}

		dbGorm.LogMode(true)
		dbGorm.DB().SetConnMaxLifetime(time.Duration(dbConf.MaxConnLifeTime) * time.Second)
		dbGorm.DB().SetMaxOpenConns(dbConf.MaxOpenConn)
		dbGorm.DB().SetMaxIdleConns(dbConf.MaxIdleConn)

		DBMapPool[confName] = dbpool
		GORMMapPool[confName] = dbGorm
	}

	if dbPool, err := GetDBPool("defualt"); err == nil {
		DBDefaultPool = dbPool
	}
	if dbPool, err := GetGormPool("defualt"); err == nil {
		GORMDefaultPool = dbPool
	}

	return nil
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbPool, ok := DBMapPool[name]; ok {
		return dbPool, nil
	}
	return nil, errors.New("get database pool error")
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbPool, ok := GORMMapPool[name]; ok {
		return dbPool, nil
	}
	return nil, errors.New("get gorm pool error")
}

func CloseDB() error {
	for _, dbPool := range DBMapPool {
		err := dbPool.Close()
		if err != nil {
			return err
		}
	}
	for _, dbPool := range GORMMapPool {
		err := dbPool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func DBPoolLogQuery(trace *TraceContext, sqlDb *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	startExecTime := time.Now()
	rows, err := sqlDb.Query(query, args...)
	endExecTime := time.Now()
	if err != nil {
		Log.TagError(trace, "mysql_success", map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		Log.TagInfo(trace, "mysql_success", map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return rows, err
}
