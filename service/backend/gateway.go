package backend

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizationGateway(accessToken string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if accessToken != "" {
			token := ctx.Request.Header.Get("Authorization")
			successToken := fmt.Sprintf("Bearer %s", accessToken)
			if token != successToken {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			}
			ctx.Next()
		}
	}
}
