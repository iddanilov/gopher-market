package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"

	"github.com/gopher-market/internal/app"
	"github.com/gopher-market/internal/config"
	"github.com/gopher-market/pkg/logging"
)

func main() {
	ctx := context.Background()
	//ctx, cancel := context.WithCancel(context.Background())
	//ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
	//defer cancel()

	log.Println("config initialing")
	cfg := config.GetConfig()

	log.Println("logger initialing")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(ctx, cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Running application   ")
	a.Run()
}
