package main

import (
	"fmt"

	"github.com/Deve-Lite/DashboardX-API-PoC/config"
	_ "github.com/Deve-Lite/DashboardX-API-PoC/docs"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/handler"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/middleware"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/postgres"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@Title		DashboardX API
//	@Version	1.0

//	@Contact.Name	Deve-Lite

//	@Host		localhost:3000
//	@BasePath	/api/v1

// @SecurityDefinitions.apikey	BearerAuth
// @In							header
// @Name						Authorization
func main() {
	cfg := config.NewConfig(".env")

	db := postgres.NewDB(cfg)
	defer db.Close()

	ch := redis.NewDB(cfg)
	defer ch.Close()

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin := gin.Default()

	app := application.NewApplication(cfg, db, ch)

	mRule := middleware.NewRule(app.AuthSrv, app.UserSrv)

	userHnd := handler.NewUserHandler(app.UserSrv, app.UserMap)
	brokerHnd := handler.NewBrokerHandler(app.BrokerSrv, app.BrokerMap)
	deviceHnd := handler.NewDeviceHandler(app.DeviceSrv, app.ControlSrv, app.DeviceMap, app.ControlMap)

	rest.NewRouter(gin, mRule, userHnd, brokerHnd, deviceHnd)

	if cfg.Server.Env == "development" {
		gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	gin.Use(middleware.CORS)

	gin.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
