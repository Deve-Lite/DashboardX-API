package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Writer.Header().Set("Content-Type", "application/json")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusOK)
		return
	}

	ctx.Next()
}
