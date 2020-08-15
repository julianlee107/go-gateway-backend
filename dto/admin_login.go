package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/julianlee107/go-gateway-backend/public"
)

type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"用户名" example:"小名" validate:"required"`
	Password string `json:"password" form:"passwd" comment:"密码" example:"123456" validate:"required"`
}

func (param *AdminLoginInput) BindValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
