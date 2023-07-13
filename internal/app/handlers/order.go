package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/dto"
	"net/http"
)

var ErrOrderAlreadyUploaded = errors.New("order already uploaded by another user")
var ErrOrderAlreadyUploadedByUser = errors.New("order already uploaded by this user")
var ErrOrderHasWrongFormat = errors.New("order already uploaded by this user")

func (h *Handler) ProcessUserOrder(c RequestContext) {
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Infof("error while reading request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading request"})
		return
	}
	orderNumber := string(requestBytes)
	userID := c.MustGet("userID").(uint)
	order, err := h.orderService.SaveOrder(dto.OrderDTO{Number: orderNumber, UserID: userID})
	if err != nil {
		if errors.Is(err, ErrOrderAlreadyUploadedByUser) {
			logger.Log.Infof("Order already uploaded by user %v", userID)
			c.JSON(http.StatusOK, gin.H{"error": "order already uploaded by this user"})
			return
		}
		if errors.Is(err, ErrOrderAlreadyUploaded) {
			logger.Log.Infof("Order already uploaded by another user")
			c.JSON(http.StatusOK, gin.H{"error": "order already uploaded by another user"})
			return
		}
		if errors.Is(err, ErrOrderHasWrongFormat) {
			logger.Log.Infof("Wrong order number format: %s", order.Number)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "wrong order number format"})
			return
		}
		logger.Log.Infof("Internal Server Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"processed": order.Number})
}

func (h *Handler) GetAllOrders(c RequestContext) {
	userID := c.MustGet("userID").(uint)
	orders, err := h.orderService.GetAllOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"error": "orders not found"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
