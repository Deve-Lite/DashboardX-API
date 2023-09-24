package application

import (
	"github.com/Deve-Lite/DashboardX-API-PoC/config"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application/mapper"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/infrastructure/cache"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/infrastructure/persistance"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/validate"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	AuthSrv    RESTAuthService
	UserSrv    UserService
	BrokerSrv  BrokerService
	DeviceSrv  DeviceService
	ControlSrv DeviceControlService

	UserMap    mapper.UserMapper
	BrokerMap  mapper.BrokerMapper
	DeviceMap  mapper.DeviceMapper
	ControlMap mapper.DeviceControlMapper
}

func NewApplication(c *config.Config, d *sqlx.DB, ch *redis.Client) *Application {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("nullmin", validate.NullMin)
		v.RegisterValidation("nullemail", validate.NullEmail)
		v.RegisterValidation("nulluuid", validate.NullUUID)
		v.RegisterValidation("nullhexcolor", validate.NullHexColor)
		v.RegisterValidation("control_attributes", validate.ControlAttributes)
		v.RegisterValidation("control_type", validate.ControlType)
		v.RegisterValidation("qos_level", validate.QoSLevel)
	}

	userRepo := persistance.NewUserRepository(d)
	brokerRepo := persistance.NewBrokerRepository(d)
	deviceRepo := persistance.NewDeviceRepository(d)
	controlRepo := persistance.NewDeviceControlRepository(d)
	tokenRepo := cache.NewTokenRepository(ch)

	authSrv := NewRESTAuthService(c, tokenRepo)
	userSrv := NewUserService(c, userRepo, authSrv)
	brokerSrv := NewBrokerService(brokerRepo)
	deviceSrv := NewDeviceService(deviceRepo, brokerSrv)
	controlSrv := NewDeviceControlService(controlRepo, deviceSrv)

	userMap := mapper.NewUserMapper()
	brokerMap := mapper.NewBrokerMapper()
	deviceMap := mapper.NewDeviceMapper()
	controlMap := mapper.NewDeviceControlMapper()

	return &Application{
		authSrv,
		userSrv,
		brokerSrv,
		deviceSrv,
		controlSrv,
		userMap,
		brokerMap,
		deviceMap,
		controlMap,
	}
}
