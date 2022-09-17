package api

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/service/backend"
)

var (
	engine *gin.Engine
)

type HttpService struct {
}

func StartHttpService() {
	if engine == nil {
		engine = gin.Default()
		engine.Use(backend.AuthorizationGateway(model.GetGlobalConfig().AccessToken))
		engine.POST("/register", backend.RegisterHandler())
		engine.OPTIONS("/exit", backend.ExitHandler())
		engine.GET("/channel")
		engine.DELETE("/channel/:channel_id", backend.RemoveChannelHandler())

		go engine.Run(fmt.Sprintf("0.0.0.0:%v", model.GetGlobalConfig().ServingPort))
	} else {
		return
	}
}

func AddService(method string, relativePath string, handlers ...gin.HandlerFunc) {
	engine.Handle(method, relativePath, handlers...)
}
