package dao

import (
	"github.com/jinzhu/gorm"
)

type App struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	AppID    string `json:"app_id" gorm:"column:app_id" description:"租户ID"`
	Name     string `json:"name" gorm:"column:name" description:"租户名"`
	Secret   string `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS string `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd      int64  `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps      int64  `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	IsDelete int8   `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
	ModifyTime
}

func (a *App) TableName() string {
	return "gateway_app"
}

func (a *App) Find(tx *gorm.DB, search *App) (*App, error) {
	app := &App{}
	err := tx.Where(search).Find(app).Error
	return app, err
}

func (a *App) Save(tx *gorm.DB) error {
	if err := tx.Save(a).Error; err != nil {
		return err
	}
	return nil
}
