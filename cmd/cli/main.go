package main

import (
	"log"
	"os"

	"github.com/Deve-Lite/DashboardX-API-PoC/config"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/cli"
	"github.com/Deve-Lite/DashboardX-API-PoC/pkg/postgres"
)

func main() {
	c := config.NewConfig(".env")

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
		cli.Seed(c)
	default:
		log.Panicf("Unknow or missing argument: %s", arg)
	}
}
