package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"net/http"
)

func (h *Handler) ProcessUserOrder(c RequestContext) {
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Infof("error while reading request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading request"})
		return
	}
	orderNumber := string(requestBytes)
	userID := c.MustGet("mustGetReturn").(uint)
	order, err := h.orderService.SaveOrder(storage.NewOrderRequest{Number: orderNumber, UserID: userID})
	if err != nil {
		switch {
		case err.Error() == "order already uploaded by this user":
			c.JSON(http.StatusOK, gin.H{"error": "order already uploaded by this user"})
			return
		case err.Error() == "order already uploaded by another user":
			c.JSON(http.StatusConflict, gin.H{"error": "order already uploaded by another user"})
			return
		case err.Error() == "order has wrong format":
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "wrong order number format"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusAccepted, gin.H{"processed": order.Number})
}

func (h *Handler) GetAllOrders(c RequestContext) {
	userID := c.MustGet("mustGetReturn").(uint)
	orders, err := h.orderService.GetAllOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"error": "orders not found"})
	}
	c.JSON(http.StatusOK, orders)
}
