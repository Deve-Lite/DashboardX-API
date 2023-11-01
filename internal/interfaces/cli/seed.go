package cli

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/Deve-Lite/DashboardX-API/pkg/postgres"
	"github.com/Deve-Lite/DashboardX-API/pkg/redis"
	"github.com/google/uuid"
)

func Seed(c *config.Config) {
	db := postgres.NewDB(c)
	defer db.Close()

	ch := redis.NewDB(c)
	defer ch.Close()

	app := application.NewApplication(c, db, ch)

	ctx := context.Background()

	puid1, _ := app.UserSrv.PreCreate(ctx, &domain.CreateUser{
		Name:     "psp515",
		Email:    "psp515@wp.pl",
		Password: "Admin123!",
		IsAdmin:  true,
	})
	puid2, _ := app.UserSrv.PreCreate(ctx, &domain.CreateUser{
		Name:     "dred",
		Email:    "dred@gmail.pl",
		Password: "User123!",
	})

	uid1, _ := app.UserSrv.Create(ctx, puid1)
	uid2, _ := app.UserSrv.Create(ctx, puid2)

	bid1, _ := app.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid1,
		Name:                "Test Devices",
		Port:                8884,
		IsSSL:               true,
		ClientID:            t.NewString("123", false, true),
		KeepAlive:           60,
		IconName:            "Home",
		IconBackgroundColor: "#ff00ff",
		Server:              "broker.hivemq.com",
	})
	app.BrokerSrv.SetCredentials(ctx, &domain.UpdateBroker{
		UserID:   uid1,
		ID:       bid1,
		Username: t.NewString("test01", false, true),
		Password: t.NewString("test01", false, true),
	})
	bid2, _ := app.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid1,
		Name:                "Home Devices",
		Port:                8884,
		IsSSL:               true,
		ClientID:            t.NewString("Admin123", false, true),
		KeepAlive:           10,
		IconName:            "Home",
		IconBackgroundColor: "#aa00ff",
		Server:              "ef57f832f11b4e89960ef452f56e6aa3.s2.eu.hivemq.cloud",
	})
	app.BrokerSrv.SetCredentials(ctx, &domain.UpdateBroker{
		UserID:   uid1,
		ID:       bid2,
		Username: t.NewString("admin", false, true),
		Password: t.NewString("Admin123", false, true),
	})
	bid3, _ := app.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid2,
		Name:                "Test Devices",
		Port:                8884,
		IsSSL:               true,
		ClientID:            t.NewString("123", false, true),
		KeepAlive:           60,
		IconName:            "Home",
		IconBackgroundColor: "#ff00ff",
		Server:              "broker.hivemq.com",
	})
	app.BrokerSrv.SetCredentials(ctx, &domain.UpdateBroker{
		UserID:   uid2,
		ID:       bid3,
		Username: t.NewString("test01", false, true),
		Password: t.NewString("test01", false, true),
	})
	bid4, _ := app.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid2,
		Name:                "Home Devices",
		Port:                8884,
		IsSSL:               true,
		ClientID:            t.NewString("Admin123", false, true),
		KeepAlive:           10,
		IconName:            "Home",
		IconBackgroundColor: "#aa00ff",
		Server:              "ef57f832f11b4e89960ef452f56e6aa3.s2.eu.hivemq.cloud",
	})
	app.BrokerSrv.SetCredentials(ctx, &domain.UpdateBroker{
		UserID:   uid2,
		ID:       bid4,
		Username: t.NewString("admin", false, true),
		Password: t.NewString("Admin123", false, true),
	})

	did1, _ := app.DeviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid1,
		Name:                "Prawdziwa lampka",
		Placing:             t.NewString("Office", false, true),
		IconName:            "Bussiness",
		IconBackgroundColor: "#86b049",
		BasePath:            t.NewString("", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid2, Valid: true},
	})
	app.DeviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid1,
		Name:                "Lamp",
		Placing:             t.NewString("Bedroom", false, true),
		IconName:            "Bussiness",
		IconBackgroundColor: "#dff5ce",
		BasePath:            t.NewString("bedroom-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid2, Valid: true},
	})
	app.DeviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid1,
		Name:                "Car",
		Placing:             t.NewString("Office", false, true),
		IconName:            "Call",
		IconBackgroundColor: "#86b049",
		BasePath:            t.NewString("car", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid1, Valid: true},
	})

	sca := make(domain.ControlAttributes)

	sca["onPayload"] = `{"state": 1}`
	sca["offPayload"] = `{"state": 0}`

	//State
	app.ControlSrv.Create(ctx, uid1, &domain.CreateDeviceControl{
		DeviceID:               did1,
		Type:                   enum.ControlState,
		Topic:                  "state",
		Name:                   "Lamp",
		QoS:                    enum.QoSZero,
		IsConfirmationRequired: false,
		IsAvailable:            true,
		IconName:               "Home",
		IconBackgroundColor:    "#aa00ff",
		CanNotifyOnPublish:     false,
		CanDisplayName:         true,
		Attributes:             sca,
	})

	//radio - mode

	lcca := make(domain.ControlAttributes)

	lcca["payloads"] = `{"White": "{"mode": 1,"data":{"brightness": 1,"extend":{"color":[256, 256, 256]}}}","RGB": "{"mode": 5,"data":{"brightness": 0.5,"extend":{"wait": 5}}}","Fade": "{"mode": 7,"data":{"brightness": 0.5,"extend":{"wait":[[255, 0, 0],[255, 255, 0],[0, 255, 0],[0, 255, 255],[0, 0, 255],[255, 0, 255]]}}}"}`

	app.ControlSrv.Create(ctx, uid1, &domain.CreateDeviceControl{
		DeviceID:               did1,
		Type:                   enum.ControlRadio,
		Topic:                  "new_mode",
		Name:                   "Lamp Mode",
		QoS:                    enum.QoSZero,
		IsConfirmationRequired: false,
		IsAvailable:            true,
		IconName:               "Home",
		IconBackgroundColor:    "#aa00ff",
		CanNotifyOnPublish:     false,
		CanDisplayName:         true,
		Attributes:             lcca,
	})

	//slider - brightness

	bca := make(domain.ControlAttributes)

	bca["minValue"] = "0"
	bca["maxValue"] = "1"
	bca["payloadTemplate"] = `{"brightness": $value}`

	app.ControlSrv.Create(ctx, uid1, &domain.CreateDeviceControl{
		DeviceID:               did1,
		Type:                   enum.ControlButton,
		Topic:                  "update_mode",
		Name:                   "Brightness",
		QoS:                    enum.QoSZero,
		IsConfirmationRequired: true,
		IsAvailable:            false,
		IconName:               "Home",
		IconBackgroundColor:    "#aa00ff",
		CanNotifyOnPublish:     false,
		CanDisplayName:         true,
		Attributes:             bca,
	})

	did2, _ := app.DeviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid2,
		Name:                "Lamp",
		Placing:             t.NewString("Office", false, true),
		IconName:            "LaptopMac",
		IconBackgroundColor: "#525b88",
		BasePath:            t.NewString("office-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid4, Valid: true},
	})
	app.DeviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid2,
		Name:                "Lamp",
		Placing:             t.NewString("Bedroom", false, true),
		IconName:            "LaptopChromebook",
		IconBackgroundColor: "#dff5ce",
		BasePath:            t.NewString("bedroom-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid4, Valid: true},
	})
	app.DeviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid2,
		Name:                "Car",
		Placing:             t.NewString("Office", false, true),
		IconName:            "LaptopChromebook",
		IconBackgroundColor: "#86B049",
		BasePath:            t.NewString("car", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid3, Valid: true},
	})

	bca2 := make(domain.ControlAttributes)

	bca2["payload"] = "No siema"

	app.ControlSrv.Create(ctx, uid2, &domain.CreateDeviceControl{
		DeviceID:               did2,
		Type:                   enum.ControlButton,
		Topic:                  "button/topic",
		Name:                   "Lamp",
		QoS:                    enum.QoSZero,
		IsConfirmationRequired: false,
		IsAvailable:            true,
		IconName:               "Home",
		IconBackgroundColor:    "#aa00ff",
		CanNotifyOnPublish:     false,
		CanDisplayName:         true,
		Attributes:             bca2,
	})
}
