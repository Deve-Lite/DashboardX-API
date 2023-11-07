package main

import (
	"flag"
	"log"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/docs"
	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest/handler"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest/middleware"
	"github.com/Deve-Lite/DashboardX-API/pkg/postgres"
	"github.com/Deve-Lite/DashboardX-API/pkg/redis"
	"github.com/Deve-Lite/DashboardX-API/pkg/smtp"
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
	cfgPath := flag.String("cfg", config.GetDefaultPath(".env"), "override default config path")
	flag.Parse()
	cfg := config.NewConfig(*cfgPath)

	db := postgres.NewDB(cfg.Postgres)
	defer db.Close()

	ch := redis.NewDB(cfg.Redis)
	defer ch.Close()

	s := smtp.NewClient(cfg.SMTP)

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin := gin.Default()

	app := application.NewApplication(cfg, db, ch, s)

	mRule := middleware.NewRule(app.AuthSrv, app.UserSrv)
	mInfo := middleware.NewInfo(cfg)

	userHnd := handler.NewUserHandler(cfg, app.UserSrv, app.UserMap)
	brokerHnd := handler.NewBrokerHandler(app.BrokerSrv, app.BrokerMap)
	deviceHnd := handler.NewDeviceHandler(app.DeviceSrv, app.ControlSrv, app.DeviceMap, app.ControlMap)

	gin.Use(middleware.CORS(cfg.CORS))

	rest.NewRouter(gin, mRule, mInfo, userHnd, brokerHnd, deviceHnd)

	setupSwagger(gin, cfg.Server)

	runServer(gin, cfg.Server)
}

func runServer(gin *gin.Engine, cfg *config.ServerConfig) {
	var err error
	if cfg.IsTLS() {
		err = gin.RunTLS(cfg.URL(), cfg.TLSCert, cfg.TLSKey)
	} else {
		err = gin.Run(cfg.URL())
	}

	if err != nil {
		log.Panic(err)
	}
}

func setupSwagger(gin *gin.Engine, cfg *config.ServerConfig) {
	if cfg.Env == "development" {
		var host string
		if cfg.DocsURLOverride != "" {
			host = cfg.DocsURLOverride
		} else {
			host = cfg.URL()
		}
		docs.SwaggerInfo.Host = host

		gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
