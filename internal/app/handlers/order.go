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
	userId := c.MustGet("userID").(int)
	order, err := h.orderService.SaveOrder(storage.NewOrderRequest{Number: orderNumber, UserID: userId})
	if err != nil {
		return
	}
	c.JSON(http.StatusAccepted, fmt.Sprintf("Accepted order: %v", order))
}
