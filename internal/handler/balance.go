package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gopher-market/internal/models"
)

func (h *Handler) getUserBalance(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}
	balance, err := h.services.Balance.GetBalance(context.Background(), strconv.Itoa(userID))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *Handler) withdraw(c *gin.Context) {
	withdrawals := models.Withdrawals{}
	userID, err := getUserID(c)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}

	responseData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal(responseData, &withdrawals)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	withdrawals.ID = strconv.Itoa(userID)

	err = h.services.Balance.Withdraw(context.Background(), withdrawals)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

	}
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) getWithdrawals(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}

	orders, err := h.services.Balance.GetWithdrawals(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}
