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

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	balance, err := h.services.Balance.GetBalance(context.Background(), strconv.Itoa(userId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (h *Handler) withdraw(c *gin.Context) {
	withdrawals := models.Withdrawals{}
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	withdrawals.Id = strconv.Itoa(userId)

	err = h.services.Balance.Withdraw(context.Background(), withdrawals)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

	}
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) getWithdrawals(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	orders, err := h.services.Balance.GetWithdrawals(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}
