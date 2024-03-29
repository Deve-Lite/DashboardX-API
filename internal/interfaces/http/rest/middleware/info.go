package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Info interface {
	Depricated(deprecation time.Time, sunset time.Time, link string) func(*gin.Context)
	Disabled(*gin.Context)
}

type info struct {
	c *config.Config
}

func NewInfo(c *config.Config) Info {
	return &info{c}
}

func (i *info) Depricated(deprecation time.Time, sunset time.Time, link string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Deprication", deprecation.UTC().Format(http.TimeFormat))
		ctx.Writer.Header().Set("Link", fmt.Sprintf(`%s%s; rel="alternate"`, i.c.Server.URL(), link))
		ctx.Writer.Header().Set("Sunset", sunset.UTC().Format(http.TimeFormat))
		ctx.Next()
	}
}

func (i *info) Disabled(ctx *gin.Context) {
	if i.c.Server.Env != "production" {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, ae.NewHTTPError(ae.ErrEndpointDisabled))
}
