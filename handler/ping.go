package handler

import (
	"template_project/model"
	"template_project/utils/render"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	UUID string `json:"uuid"`
}

func Ping(c *gin.Context) {
	fmt.Println("------------ping-----pong------")
	render.RespJson(c, 0, "success", "----template_project----pong")
}

func TestPost(ctx *gin.Context)  {
	var body interface{}
	if err := ctx.ShouldBindJSON( &body); err != nil {
		fmt.Println(err)
		render.RespJson(ctx, http.StatusInternalServerError, err.Error(), "pong")
		return
	}
	fmt.Println("post body:", body)
	render.RespJson(ctx, http.StatusOK, "ok", body)
}

func GetUser(ctx *gin.Context)  {
	dao := &model.User{}
	user, err := dao.QueryUserById(1)
	if err != nil {
		fmt.Println("-------GetUser---:err", err)
		render.RespJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	render.RespJson(ctx, http.StatusOK, "ok", user)
}
