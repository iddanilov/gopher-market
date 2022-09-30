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
			user.GET("/balance", h.getUserBalance)
			user.GET("/balance/withdraw", h.getBalanceWithdraw)
		}

		orders := api.Group("/orders")
		{
			orders.GET("/", h.GetOrders)
			orders.PUT("/", h.SaveOrder)
		}
	}

	return router
}
