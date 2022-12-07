package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	_ "github.com/gopher-market/doc"
	"github.com/gopher-market/internal/config"
	"github.com/gopher-market/internal/handlers"
	"github.com/gopher-market/internal/service"
	"github.com/gopher-market/internal/storage"
	"github.com/gopher-market/internal/storage/postgres"
	"github.com/gopher-market/pkg/logging"
)

type App struct {
	ctx        context.Context
	cfg        *config.Config
	logger     *logging.Logger
	router     *gin.Engine
	httpServer *http.Server
}

func NewApp(ctx context.Context, cfg *config.Config, logger *logging.Logger) (App, error) {
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	// migrations
	m := postgres.NewMigrationsPostgres(db)
	err = m.CreateUserTable(ctx)
	if err != nil {
		logger.Error(err)
	}
	err = m.CreateOrdersTable(ctx)
	if err != nil {
		logger.Error(err)
	}
	err = m.CreateBalanceTable(ctx)
	if err != nil {
		logger.Error(err)
	}
	err = m.CreateWithdrawalsTable(ctx)

	if err != nil {
		logger.Error(err)
	}

	logger.Info("router initializing")
	newStorage := storage.NewStorage(ctx, db, logger)
	services := service.NewService(ctx, newStorage, cfg, logger)
	routers := handlers.NewHandler(ctx, services, logger)

	return App{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
		router: routers.InitRoutes(),
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	a.logger.Infof("bind application to host and port: %s", a.cfg.Listen.RunAddress)
	var err error
	listener, err = net.Listen("tcp", a.cfg.Listen.RunAddress)
	if err != nil {
		a.logger.Info(err)
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080", fmt.Sprintf("http://%s", a.cfg.Listen.RunAddress)},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err = a.httpServer.Shutdown(a.ctx)
	if err != nil {
		a.logger.Fatal(err)
	}
}
