package main

import (
	"fmt"

	"github.com/Deve-Lite/DashboardX-API-PoC/config"
	_ "github.com/Deve-Lite/DashboardX-API-PoC/docs"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/infrastructure/persistance"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/auth"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/handler"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/middleware"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/mapper"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/postgres"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
	cfg := config.NewConfig()

	db := postgres.NewDB(cfg)
	defer db.Close()

	gin := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("nullmin", validate.NullMin)
		v.RegisterValidation("nullemail", validate.NullEmail)
		v.RegisterValidation("nulluuid", validate.NullUUID)
		v.RegisterValidation("nullhexcolor", validate.NullHexColor)
		v.RegisterValidation("control_attributes", validate.ControlAttributes)
		v.RegisterValidation("control_type", validate.ControlType)
		v.RegisterValidation("qos_level", validate.QoSLevel)
	}

	auth := auth.NewRESTAuth(cfg)

	userRepo := persistance.NewUserRepository(db)
	brokerRepo := persistance.NewBrokerRepository(db)
	deviceRepo := persistance.NewDeviceRepository(db)
	controlRepo := persistance.NewDeviceControlRepository(db)

	userSrv := application.NewUserService(cfg, userRepo, auth)
	brokerSrv := application.NewBrokerService(brokerRepo)
	deviceSrv := application.NewDeviceService(deviceRepo, brokerRepo)
	controlSrv := application.NewDeviceControlService(deviceRepo, controlRepo)

	userMap := mapper.NewUserMapper()
	brokerMap := mapper.NewBrokerMapper()
	deviceMap := mapper.NewDeviceMapper()
	controlMap := mapper.NewDeviceControlMapper()

	userHnd := handler.NewUserHandler(userSrv, userMap)
	brokerHnd := handler.NewBrokerHandler(brokerSrv, brokerMap)
	deviceHnd := handler.NewDeviceHandler(deviceSrv, controlSrv, deviceMap, controlMap)

	mRule := middleware.NewRule(auth, userSrv)

	rest.NewRouter(gin, mRule, userHnd, brokerHnd, deviceHnd)

	if cfg.Server.Env == "development" {
		gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	gin.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
