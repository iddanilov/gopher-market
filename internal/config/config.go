package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE" env-default:"port" env-description:"port or sock"`
		BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"PORT" env-default:"10000"`
		SocketFile string `env:"SOCKET_FILE " env-default:"app.sock"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL"`
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-required:"true"`
			Password string `env:"ADMIN_PWD" env-required:"true"`
		}
	}
	Postgres struct {
		DSN string `env:"POSTGRES_DSN"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("config")
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "Monolith system"
			description, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(description)
			log.Fatal(err)
		}
	})
	return instance
}
