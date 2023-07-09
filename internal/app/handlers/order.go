package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
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
	c.JSON(http.StatusAccepted, "Accepted request: "+orderNumber)
}
