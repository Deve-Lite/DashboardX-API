package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Server      *ServerConfig
	Postgres    *PostgresConfig
	Crytpo      *CryptoConfig
	JWT         *JWTConfig
	Redis       *RedisConfig
	SMTP        *SMTPConfig
	Frontend    *FrontendConfig
	MailAddress *MailAddressConfig
	CORS        *CORSConfig
}

type ServerConfig struct {
	Host            string `mapstructure:"HOST"`
	Port            uint16 `mapstructure:"PORT"`
	Domain          string `mapstructure:"DOMAIN"`
	Env             string `mapstructure:"ENV"`
	TLSCert         string `mapstructure:"TLS_CERT"`
	TLSKey          string `mapstructure:"TLS_KEY"`
	DocsURLOverride string `mapstructure:"DOCS_URL_OVERRIDE"`
}

func (c *ServerConfig) URL() string {
	url := c.Host
	if c.Port != 0 && c.Port != 80 && c.Port != 443 {
		url += fmt.Sprintf(":%d", c.Port)
	}
	return url
}

func (c *ServerConfig) IsTLS() bool {
	if _, err := os.Stat(c.TLSCert); err != nil {
		return false
	}
	if _, err := os.Stat(c.TLSKey); err != nil {
		return false
	}
	return true
}

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     uint16 `mapstructure:"POSTGRES_PORT"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	User     string `mapstructure:"POSTGRES_USER"`
	Database string `mapstructure:"POSTGRES_DATABASE"`
}

type CryptoConfig struct {
	HashCost      uint8  `mapstructure:"CRYPTO_HASH_COST"`
	BrokersAESKey string `mapstructure:"CRYPTO_BROKERS_AES_KEY"`
}

type JWTConfig struct {
	AccessSecret         string  `mapstructure:"JWT_ACCESS_SECRET"`
	AccessLifespanHours  float32 `mapstructure:"JWT_ACCESS_LIFESPAN_HOURS"`
	RefreshSecret        string  `mapstructure:"JWT_REFRESH_SECRET"`
	RefreshLifespanHours float32 `mapstructure:"JWT_REFRESH_LIFESPAN_HOURS"`
	ConfirmSecret        string  `mapstructure:"JWT_CONFIRM_SECRET"`
	ConfirmLifespanHours float32 `mapstructure:"JWT_CONFIRM_LIFESPAN_HOURS"`
	ResetSecret          string  `mapstructure:"JWT_RESET_SECRET"`
	ResetLifespanMinutes float32 `mapstructure:"JWT_RESET_LIFESPAN_MINUTES"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     uint16 `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	Database uint8  `mapstructure:"REDIS_DATABASE"`
}

type SMTPConfig struct {
	User               string `mapstructure:"SMTP_USER"`
	Password           string `mapstructure:"SMTP_PASSWORD"`
	Port               uint16 `mapstructure:"SMTP_PORT"`
	Host               string `mapstructure:"SMTP_HOST"`
	InsecureSkipVerify bool   `mapstructure:"SMTP_INSECURE_SKIP_VERIFY"`
}

type FrontendConfig struct {
	ConfirmAccountURL string `mapstructure:"FRONTEND_CONFIRM_ACCOUNT_URL"`
	ResetPasswordURL  string `mapstructure:"FRONTEND_RESET_PASSWORD_URL"`
}

type MailAddressConfig struct {
	Default string `mapstructure:"MAIL_ADDRESS_DEFAULT"`
}

type CORSConfig struct {
	Credentials string `mapstructure:"CORS_CREDENTIALS"`
	Methods     string `mapstructure:"CORS_METHODS"`
	Origin      string `mapstructure:"CORS_ORIGIN"`
	Headers     string `mapstructure:"CORS_HEADERS"`
}

func loadConfig[T interface{}](v *viper.Viper, c T) *T {
	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatal("Config can not be loaded, Error: ", err)
	}

	return &c
}

func NewConfig(filename string) *Config {
	v := viper.GetViper()

	_, b, _, _ := runtime.Caller(0)
	v.SetConfigFile(filepath.Join(filepath.Dir(b), "..", filename))

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Can not find the .env file. Error: ", err)
	}

	config := Config{
		Server:      loadConfig(v, ServerConfig{}),
		Postgres:    loadConfig(v, PostgresConfig{}),
		Crytpo:      loadConfig(v, CryptoConfig{}),
		JWT:         loadConfig(v, JWTConfig{}),
		Redis:       loadConfig(v, RedisConfig{}),
		SMTP:        loadConfig(v, SMTPConfig{}),
		Frontend:    loadConfig(v, FrontendConfig{}),
		MailAddress: loadConfig(v, MailAddressConfig{}),
		CORS:        loadConfig(v, CORSConfig{}),
	}

	return &config
}
