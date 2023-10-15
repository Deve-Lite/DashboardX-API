package postgres

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/golang-migrate/migrate"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	DuplicatedKey = "23505"
)

const (
	UserEmailConstraint          = "users_email_key"
	BrokerUserIDServerConstraint = "brokers_user_id_server_key"
)

func NewDB(c *config.Config) *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Database)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}

	log.Print("Postgres connected successfully")

	return db
}

func RunUp(c *config.Config) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database)

	m, err := migrate.New(getMigrationsPath(), url)
	if err != nil {
		log.Fatalf("Migrate: connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Print("Migrate: no change")
		return
	}

	log.Print("Migrate: up success")
}

func RunDown(c *config.Config) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database)

	m, err := migrate.New(getMigrationsPath(), url)
	if err != nil {
		log.Fatalf("Migrate: connect error: %s", err)
	}

	err = m.Down()
	defer m.Close()
	if err != nil {
		log.Fatalf("Migrate: down error: %s", err)
	}

	log.Print("Migrate: down success")
}

func Create(c *config.Config) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, c.Postgres.Database))
	if err != nil {
		log.Panic("Could not create database. Error: ", err)
	}

	log.Printf(`Database "%s" created`, c.Postgres.Database)
}

func Drop(c *config.Config) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, c.Postgres.Database))
	if err != nil {
		log.Panic("Could not drop database. Error: ", err)
	}

	log.Printf(`Database "%s" dropped`, c.Postgres.Database)
}

func getMigrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	p := filepath.Join(filepath.Dir(b), "../..", "migrations")
	return "file://" + strings.ReplaceAll(p, "\\", "/")
}
