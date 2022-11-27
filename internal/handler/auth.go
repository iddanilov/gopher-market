package handler

import (
	"fmt"
	"github.com/gopher-market/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registration(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		//newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		if err.Error() == "user already registered" {
			authUserAlreadyRegisteredErrorResponse(c, err.Error())
		} else {
			newErrorResponse(c, http.StatusOK, err.Error())
		}
		return
	}

	token, err := h.getAuthToken(input.Login, input.Password)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}
	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) login(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		//newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.getAuthToken(input.Login, input.Password)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) getAuthToken(login, password string) (string, error) {
	return h.services.Authorization.GenerateToken(login, password)
}
