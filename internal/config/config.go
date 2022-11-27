package config

import (
	goflag "flag"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	flag "github.com/spf13/pflag"
)

var RunAddress = flag.StringP("a", "a", "localhost:8080", "help RunAddress for flagname")
var DatabaseURI = flag.StringP("d", "d", "", "help db for flagname")
var Accrual = flag.StringP("r", "r", "", "help message for flagname")

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE" env-default:"port" env-description:"port or sock"`
		RunAddress string `env:"RUN_ADDRESS" env-default:""`
		SocketFile string `env:"SOCKET_FILE " env-default:"app.sock"`
	}
	AppConfig struct {
		LogLevel string `env:"LOG_LEVEL" env-default:"debug"`
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
		flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
		flag.Parse()
	})

	if instance.Listen.RunAddress == "" {
		instance.Listen.RunAddress = *RunAddress
	}
	if instance.Postgres.DSN == "" {
		instance.Postgres.DSN = *DatabaseURI
	}

	if instance.Accrual.Address == "" {
		instance.Accrual.Address = *Accrual
	}
	return instance
}
