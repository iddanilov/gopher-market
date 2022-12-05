package handler

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (h *Handler) loadOrder(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}
	responseData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	orderInt, err := strconv.Atoi(string(responseData))
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, "")
	}
	if !valid(orderInt) {
		newErrorResponse(c, http.StatusUnprocessableEntity, "")
		return
	}
	err = h.services.Orders.LoadOrder(userID, string(responseData))
	if err != nil {
		if err.Error() == "order already loaded" {
			c.JSON(http.StatusConflict, nil)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) getOrders(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		authErrorResponse(c, err.Error())
		return
	}

	orders, err := h.services.Orders.GetOrders(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(*orders) <= 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	c.JSON(http.StatusOK, orders)

}

func valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
