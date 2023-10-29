package rest

import (
	"time"

	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest/handler"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	g *gin.Engine,
	mr middleware.Rule,
	mi middleware.Info,
	uh handler.UserHandler,
	bh handler.BrokerHandler,
	dh handler.DeviceHandler) {
	r := g.Group("/api/v1")

	// User API
	ug := r.Group("users")
	ug.POST("/register", uh.Register)
	ug.POST("/login", uh.Login)
	ug.POST("/confirm-account", mr.ValidConfirm, uh.ConfirmAccount)
	ug.POST("/confirm-account/resend", uh.ResendConfirmAccount)
	ug.POST("/refresh", mi.Depricated(
		time.Date(2023, 10, 26, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 11, 26, 23, 59, 59, 0, time.UTC),
		"/api/v1/users/me/tokens"), mr.ValidRefresh, uh.Tokens)
	ug.POST("/reset-password", uh.ResetPasswordToken)
	ug.PATCH("/reset-password", mr.ValidReset, mr.ValidResetSubject, uh.ResetPasswordChange)
	ug.GET("/me", mr.LoggedIn, uh.Get)
	ug.PATCH("/me", mr.LoggedIn, uh.Update)
	ug.DELETE("/me", mr.LoggedIn, uh.Delete)
	ug.POST("/me/tokens", mr.ValidRefresh, uh.Tokens)
	ug.PATCH("/me/password", mr.LoggedIn, uh.ChangePassword)

	// Broker API
	bg := r.Group("brokers")
	bg.GET("", mr.LoggedIn, bh.List)
	bg.POST("", mr.LoggedIn, bh.Create)
	bg.GET("/:brokerId", mr.LoggedIn, bh.Get)
	bg.PATCH("/:brokerId", mr.LoggedIn, bh.Update)
	bg.DELETE("/:brokerId", mr.LoggedIn, bh.Delete)
	bg.GET("/:brokerId/credentials", mr.LoggedIn, bh.GetCredentials)
	bg.PUT("/:brokerId/credentials", mr.LoggedIn, bh.SetCredentials)

	// Device API
	dg := r.Group("devices")
	dg.GET("", mr.LoggedIn, dh.List)
	dg.POST("", mr.LoggedIn, dh.Create)
	dg.GET("/:deviceId", mr.LoggedIn, dh.Get)
	dg.PATCH("/:deviceId", mr.LoggedIn, dh.Update)
	dg.DELETE("/:deviceId", mr.LoggedIn, dh.Delete)
	dg.GET("/:deviceId/controls", mr.LoggedIn, dh.ListControls)
	dg.POST("/:deviceId/controls", mr.LoggedIn, dh.CreateControl)
	dg.PATCH("/:deviceId/controls/:controlId", mr.LoggedIn, dh.UpdateControl)
	dg.DELETE("/:deviceId/controls/:controlId", mr.LoggedIn, dh.DeleteControl)
}
