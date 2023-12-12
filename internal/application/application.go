package application

import (
	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/mapper"
	"github.com/Deve-Lite/DashboardX-API/internal/infrastructure/cache"
	"github.com/Deve-Lite/DashboardX-API/internal/infrastructure/persistance"
	ismtp "github.com/Deve-Lite/DashboardX-API/internal/infrastructure/smtp"
	"github.com/Deve-Lite/DashboardX-API/pkg/smtp"
	"github.com/Deve-Lite/DashboardX-API/pkg/validate"
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
	EventSrv   EventService

	UserMap    mapper.UserMapper
	BrokerMap  mapper.BrokerMapper
	DeviceMap  mapper.DeviceMapper
	ControlMap mapper.DeviceControlMapper
}

func NewApplication(c *config.Config, d *sqlx.DB, ch *redis.Client, s smtp.Client) *Application {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("emptymin", validate.EmptyMin)
		v.RegisterValidation("emptyemail", validate.EmptyEmail)
		v.RegisterValidation("emptyuuid", validate.EmptyUUID)
		v.RegisterValidation("emptyhexcolor", validate.EmptyHexColor)
		v.RegisterValidation("control_attributes", validate.ControlAttributes)
		v.RegisterValidation("control_type", validate.ControlType)
		v.RegisterValidation("qos_level", validate.QoSLevel)
		v.RegisterValidation("requirednullstring", validate.RequiredNullString)
	}

	userRepo := persistance.NewUserRepository(d)
	brokerRepo := persistance.NewBrokerRepository(d)
	deviceRepo := persistance.NewDeviceRepository(d)
	controlRepo := persistance.NewDeviceControlRepository(d)
	tokenRepo := cache.NewTokenRepository(ch)
	preUserRepo := cache.NewPreUserRepository(ch)
	userActionRepo := cache.NewUserActionRepository(ch)

	mailAdp := ismtp.NewMailAdapter(c, s)

	eventSrv := NewEventService()
	mailSrv := NewMailService(mailAdp)
	cryptoSrv := NewCryptoService(c)
	authSrv := NewRESTAuthService(c, tokenRepo, cryptoSrv)
	userSrv := NewUserService(c, preUserRepo, userRepo, userActionRepo,
		authSrv, mailSrv, cryptoSrv, eventSrv)
	brokerSrv := NewBrokerService(c, brokerRepo, cryptoSrv, eventSrv)
	deviceSrv := NewDeviceService(deviceRepo, brokerSrv, eventSrv)
	controlSrv := NewDeviceControlService(controlRepo, deviceSrv, eventSrv)

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
		eventSrv,
		userMap,
		brokerMap,
		deviceMap,
		controlMap,
	}
}
