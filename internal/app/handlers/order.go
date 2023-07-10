package handlers

import (
	"fmt"
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
	userID := c.MustGet("userID").(int)
	order, err := h.orderService.SaveOrder(storage.NewOrderRequest{Number: orderNumber, UserID: userID})
	if err != nil {
		switch {
		case err.Error() == "order already uploaded by this user":
			c.AbortWithStatus(http.StatusOK)
		case err.Error() == "order already uploaded by another user":
			c.AbortWithStatus(http.StatusConflict)
		case err.Error() == "order has wrong format":
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	c.JSON(http.StatusAccepted, fmt.Sprintf("Accepted order: %v", order))
}
