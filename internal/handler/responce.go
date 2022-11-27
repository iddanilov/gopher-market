package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func authErrorResponse(c *gin.Context, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{message})
}

func authUserAlreadyRegisteredErrorResponse(c *gin.Context, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(http.StatusConflict, errorResponse{message})
}
