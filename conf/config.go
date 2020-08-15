package conf

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg       *ini.File
	Databases *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize int

	Mysql = &Database{}
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("failed to parse app.ini: %v", err)
	}
	Databases, err = ini.Load("conf/databases.ini")
	if err != nil {
		log.Fatalf("failed to parse databases.ini: %v", err)
	}
	LoadBase()
	LoadServer()
	LoadApp()
	LoadDatabase()
}
func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("failed to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60))
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60))

}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Failed to get section 'app': %v", err)
	}
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

type Database struct {
	TYPE               string
	USER               string
	PASSWORD           string
	HOST               string
	NAME               string
	MAXIDELCONNECTIONS int
	MAXOPENCONNECTIONS int
}

func LoadDatabase() {
	mysql, err := Databases.GetSection("mysql")
	if err != nil {
		log.Fatalf("Failed to get section 'mysql': %v", err)
	}
	err = mysql.MapTo(Mysql)
	if err != nil {
		log.Fatalf("Failed to load section 'mysql': %v", err)
	}
}
