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
