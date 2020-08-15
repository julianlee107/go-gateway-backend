package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/julianlee107/go-common/lib"
	"net/http"
	"strings"
)

type ResponseCode int

type Response struct {
	Code  ResponseCode `json:"code"`
	Msg   string       `json:"msg"`
	Data  interface{}  `json:"data"`
	Stack interface{}  `json:"stack"`
}

const (
	SuccessCode ResponseCode = iota + 1000
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode
)

func ResponseError(c *gin.Context, code ResponseCode, err error) {
	stack := ""
	if c.Query("is_debug") == "1" || lib.GetConfEnv() == "debug" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}
	resp := &Response{
		Code:  code,
		Msg:   err.Error(),
		Data:  "",
		Stack: stack,
	}
	c.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(http.StatusOK, err)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	resp := &Response{
		Code:  SuccessCode,
		Msg:   "Success",
		Data:  data,
		Stack: "",
	}
	c.JSON(http.StatusOK, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
