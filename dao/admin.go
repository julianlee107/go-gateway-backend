package dao

import (
	"errors"
	"github.com/julianlee107/gatewayScaffold/dto"
	"github.com/julianlee107/gatewayScaffold/public"
)
import "github.com/jinzhu/gorm"

type Admin struct {
	Id       int    `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName string `json:"user_name" gorm:"column:user_name" description:"用户名"`
	Salt     string `json:"salt" gorm:"column:salt" description:"盐"`
	Password string `json:"password" gorm:"column:password" description:"密码"`
	IsDelete int8   `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
	ModifyTime
}

func (a *Admin) TableName() string {
	return "gateway_admin"
}

func (a *Admin) LoginCheck(tx *gorm.DB, params *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := a.Find(tx, &Admin{UserName: params.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码错误")
	}
	return adminInfo, nil
}

func (a *Admin) Find(tx *gorm.DB, query *Admin) (*Admin, error) {
	result := &Admin{}
	err := tx.Where(query).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *Admin) Save(tx *gorm.DB) error {
	return tx.Save(a).Error
}
