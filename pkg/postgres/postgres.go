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

func NewDB(c *config.PostgresConfig) *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database)

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

func RunUp(c *config.PostgresConfig) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database)

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

func RunDown(c *config.PostgresConfig) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database)

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

func Create(c *config.PostgresConfig) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, c.Database))
	if err != nil {
		log.Panic("Could not create database. Error: ", err)
	}

	log.Printf(`Database "%s" created`, c.Database)
}

func Drop(c *config.PostgresConfig) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, c.Database))
	if err != nil {
		log.Panic("Could not drop database. Error: ", err)
	}

	log.Printf(`Database "%s" dropped`, c.Database)
}

type Database struct {
	Name string `db:"datname"`
}

func Exists(c *config.PostgresConfig) bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic("Can not connect to Postgres. Error: ", err)
	}
	defer db.Close()

	d := Database{}
	db.Get(&d, fmt.Sprintf(`SELECT "datname" FROM "pg_catalog"."pg_database" WHERE lower("datname") = lower('%s')`, c.Database))

	return strings.EqualFold(d.Name, c.Database)
}

func getMigrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	p := filepath.Join(filepath.Dir(b), "../..", "migrations")
	return "file://" + strings.ReplaceAll(p, "\\", "/")
}
