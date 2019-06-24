package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespJsonData struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func RespJson(c *gin.Context, code int, msg string, data interface{}) {
	result := &RespJsonData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, result)
}

func RespJsonWithError(c *gin.Context, code int, msg string) {
	RespJson(c, code, msg, nil)
}

func RespJsonWithBindingError(c *gin.Context, code int, err error) {
	RespJson(c, code, err.Error(), nil)
}
