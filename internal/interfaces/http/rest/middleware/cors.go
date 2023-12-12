package middleware

import (
	"net/http"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.CORSConfig) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", cfg.Credentials)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", cfg.Methods)
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", cfg.Origin)
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", cfg.Headers)
		// ctx.Writer.Header().Set("Content-Type", "application/json")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
