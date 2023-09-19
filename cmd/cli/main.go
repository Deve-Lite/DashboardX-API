package main

import (
	"context"
	"log"
	"os"

	"github.com/Deve-Lite/DashboardX-API-PoC/config"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/infrastructure/persistance"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/auth"
	t "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/postgres"
	"github.com/google/uuid"
)

func main() {
	c := config.NewConfig()

	switch arg := os.Args[1]; arg {
	case "migrate":
	case "up":
		postgres.RunUp(c)
	case "rollback":
	case "down":
		postgres.RunDown(c)
	case "create":
		postgres.Create(c)
	case "seed":
		seed(c)
	default:
		log.Panicf("Unknow or missing argument: %s", arg)
	}
}

func seed(c *config.Config) {
	auth := auth.NewRESTAuth(c)

	db := postgres.NewDB(c)

	userRepo := persistance.NewUserRepository(db)
	brokerRepo := persistance.NewBrokerRepository(db)
	deviceRepo := persistance.NewDeviceRepository(db)

	userSrv := application.NewUserService(c, userRepo, auth)
	brokerSrv := application.NewBrokerService(brokerRepo)
	deviceSrv := application.NewDeviceService(deviceRepo, brokerRepo)

	ctx := context.Background()

	uid1, _ := userSrv.Create(ctx, &domain.CreateUser{
		Name:     "psp515",
		Email:    "psp515@wp.pl",
		Password: "Admin123!",
		IsAdmin:  true,
	})
	uid2, _ := userSrv.Create(ctx, &domain.CreateUser{
		Name:     "dred",
		Email:    "dred@gmail.pl",
		Password: "User123!",
	})

	bid1, _ := brokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid1,
		Name:                "Test Devices",
		Port:                8883,
		Username:            t.NewString("test01", false, true),
		IsSSL:               true,
		Password:            t.NewString("test01", false, true),
		ClientID:            t.NewString("123", false, true),
		KeepAlive:           60,
		IconName:            "default.png",
		IconBackgroundColor: "#ff00ff",
		Server:              "ef57f832f11b4e89960ef452f56e6aa3.s2.eu.hivemq.cloud",
	})
	bid2, _ := brokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid1,
		Name:                "Home Devices",
		Port:                8883,
		Username:            t.NewString("test01", false, true),
		IsSSL:               true,
		Password:            t.NewString("test01", false, true),
		ClientID:            t.NewString("123", false, true),
		KeepAlive:           10,
		IconName:            "default.png",
		IconBackgroundColor: "#aa00ff",
		Server:              "ef57f832f11b4e89960ef452f56e6aa3.s2.eu.hivemq.cloud",
	})
	bid3, _ := brokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid2,
		Name:                "Test Devices",
		Port:                8883,
		Username:            t.NewString("test01", false, true),
		IsSSL:               true,
		Password:            t.NewString("test01", false, true),
		ClientID:            t.NewString("123", false, true),
		KeepAlive:           20,
		IconName:            "default.png",
		IconBackgroundColor: "#cc00ff",
		Server:              "ef57f832f11b4e89960ef452f56e6aa3.s2.eu.hivemq.cloud",
	})
	bid4, _ := brokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              uid2,
		Name:                "Test Devices",
		Port:                8883,
		Username:            t.NewString("test01", false, true),
		IsSSL:               true,
		Password:            t.NewString("test01", false, true),
		ClientID:            t.NewString("123", false, true),
		KeepAlive:           60,
		IconName:            "default.png",
		IconBackgroundColor: "#dd00ff",
		Server:              "ef57f832f11b4e89960ef452f56e6aa3.s2.eu.hivemq.cloud",
	})

	deviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid1,
		Name:                "Lamp",
		Placing:             t.NewString("Office", false, true),
		IconName:            "default2.png",
		IconBackgroundColor: "#86b049",
		BasePath:            t.NewString("office-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid2, Valid: true},
	})
	deviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid1,
		Name:                "Lamp",
		Placing:             t.NewString("Bedroom", false, true),
		IconName:            "default2.png",
		IconBackgroundColor: "#dff5ce",
		BasePath:            t.NewString("bedroom-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid2, Valid: true},
	})
	deviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid1,
		Name:                "Car",
		Placing:             t.NewString("Office", false, true),
		IconName:            "default2.png",
		IconBackgroundColor: "#86b049",
		BasePath:            t.NewString("car", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid1, Valid: true},
	})
	deviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid2,
		Name:                "Lamp",
		Placing:             t.NewString("Office", false, true),
		IconName:            "default2.png",
		IconBackgroundColor: "#525b88",
		BasePath:            t.NewString("office-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid4, Valid: true},
	})
	deviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid2,
		Name:                "Lamp",
		Placing:             t.NewString("Bedroom", false, true),
		IconName:            "default2.png",
		IconBackgroundColor: "#dff5ce",
		BasePath:            t.NewString("bedroom-lamp", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid4, Valid: true},
	})
	deviceSrv.Create(ctx, &domain.CreateDevice{
		UserID:              uid2,
		Name:                "Car",
		Placing:             t.NewString("Office", false, true),
		IconName:            "default2.png",
		IconBackgroundColor: "#86B049",
		BasePath:            t.NewString("car", false, true),
		BrokerID:            uuid.NullUUID{UUID: bid3, Valid: true},
	})
}
