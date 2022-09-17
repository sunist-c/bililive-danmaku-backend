package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/bililive-danmaku/model"
	"time"
)

func ExitHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, nil)
		defer func() {
			time.Sleep(time.Second)
			model.WaitGroup.Done()
		}()
	}
}
