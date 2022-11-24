package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	flag "github.com/spf13/pflag"
)

var RunAddress = flag.StringP("a", "a", "127.0.0.1:10000", "help RunAddress for flagname")
var DatabaseURI = flag.StringP("d", "d", "", "help db for flagname")
var Accrual = flag.StringP("r", "r", "127.0.0.1:8080", "help message for flagname")

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE" env-default:"port" env-description:"port or sock"`
		RunAddress string `env:"RUN_ADDRESS" env-default:"0.0.0.0:10000"`
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
		DSN string `env:"DATABASE_URI"`
	}
	Accrual struct {
		Address string `env:"ACCRUAL_SYSTEM_ADDRESS"`
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

	if instance.Listen.RunAddress == "0.0.0.0:10000" || *RunAddress != "" {
		instance.Listen.RunAddress = *RunAddress
	}
	if instance.Postgres.DSN == "" || *DatabaseURI != "" {
		instance.Postgres.DSN = *DatabaseURI
	}
	if instance.Accrual.Address == "" {
		instance.Accrual.Address = *Accrual
	}
	return instance
}
