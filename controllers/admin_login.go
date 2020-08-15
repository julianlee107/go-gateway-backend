package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/julianlee107/go-gateway-backend/dto"
	"github.com/julianlee107/go-gateway-backend/middleware"
	"net/http"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("admin_login/login", adminLogin.AdminLogin)

}

// AdminLogin
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (admin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParams(c); err != nil {
		middleware.ResponseError(c, http.StatusBadRequest, err)
		return
	}
	middleware.ResponseSuccess(c, "")

}
