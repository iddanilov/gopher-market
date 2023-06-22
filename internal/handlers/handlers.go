package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/gopher-market/doc"
	"github.com/gopher-market/internal/service"
	"github.com/gopher-market/pkg/logging"
)

type Handler struct {
	ctx      context.Context
	services *service.Service
	logger   *logging.Logger
}

func NewHandler(ctx context.Context, services *service.Service, logger *logging.Logger) *Handler {
	return &Handler{
		ctx:      ctx,
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", h.registration)
			user.POST("/login", h.login)
		}

		orders := user.Group("/orders", h.userIdentity)
		{
			orders.POST("/", h.loadOrder)
			orders.GET("/", h.getOrders)
		}

		balance := user.Group("/balance", h.userIdentity)
		{
			balance.GET("/", h.getUserBalance)
			balance.POST("/withdraw", h.withdraw)
		}
		withdrawals := user.Group("/withdrawals", h.userIdentity)
		{
			withdrawals.GET("/", h.getWithdrawals)
		}
	}

	return router
}
