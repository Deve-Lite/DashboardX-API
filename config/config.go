package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   *ServerConfig
	Postgres *PostgresConfig
	Crytpo   *CryptoConfig
	JWT      *JWTConfig
	Redis    *RedisConfig
}

type ServerConfig struct {
	Host string `mapstructure:"HOST"`
	Port uint16 `mapstructure:"PORT"`
	Env  string `mapstructure:"ENV"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     uint16 `mapstructure:"POSTGRES_PORT"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	User     string `mapstructure:"POSTGRES_USER"`
	Database string `mapstructure:"POSTGRES_DATABASE"`
}

type CryptoConfig struct {
	HashCost uint8 `mapstructure:"CRYPTO_HASH_COST"`
}

type JWTConfig struct {
	AccessSecret    string  `mapstructure:"JWT_ACCESS_SECRET"`
	AccessLifespan  float32 `mapstructure:"JWT_ACCESS_LIFESPAN"`
	RefreshSecret   string  `mapstructure:"JWT_REFRESH_SECRET"`
	RefreshLifespan float32 `mapstructure:"JWT_REFRESH_LIFESPAN"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     uint16 `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

func loadConfig[T interface{}](v *viper.Viper, c T) *T {
	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatal("Config can not be loaded, Error: ", err)
	}

	return &c
}

func NewConfig() *Config {
	v := viper.GetViper()

	v.AddConfigPath(".")
	v.SetConfigFile(".env")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Can not find the .env file. Error: ", err)
	}

	config := Config{
		Server:   loadConfig(v, ServerConfig{}),
		Postgres: loadConfig(v, PostgresConfig{}),
		Crytpo:   loadConfig(v, CryptoConfig{}),
		JWT:      loadConfig(v, JWTConfig{}),
		Redis:    loadConfig(v, RedisConfig{}),
	}

	return &config
}
