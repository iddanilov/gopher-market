package handler

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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
	err = h.services.Orders.LoadOrder(userID, string(responseData))
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"orders_pkey\"" {
			c.JSON(http.StatusOK, nil)
			return
		} else if err.Error() == "order already loaded" {
			c.JSON(http.StatusConflict, nil)
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
