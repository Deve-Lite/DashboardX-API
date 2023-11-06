package main

import (
	"flag"
	"log"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/cli"
	"github.com/Deve-Lite/DashboardX-API/pkg/postgres"
)

func main() {
	cfgPath := flag.String("cfg", config.GetDefaultPath(".env"), "override default config path")
	oper := flag.String("op", "", "set an operation to run")
	flag.Parse()
	c := config.NewConfig(*cfgPath)

	switch *oper {
	case "migrate":
	case "up":
		postgres.RunUp(c.Postgres)
	case "rollback":
	case "down":
		postgres.RunDown(c.Postgres)
	case "create":
		postgres.Create(c.Postgres)
	case "seed":
		cli.Seed(c)
	default:
		log.Panicf("You have to specify an operation: %s", *oper)
	}
}
