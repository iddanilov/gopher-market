package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gopher-market/doc"
	"github.com/gopher-market/internal/service"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
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
		balance.GET("/withdrawals", h.getWithdrawals)
	}

	return router
}
