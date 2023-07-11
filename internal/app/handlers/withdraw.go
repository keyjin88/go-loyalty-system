package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"net/http"
)

func (h *Handler) SaveWithdraw(c RequestContext) {
	var req storage.WithdrawRequest
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Infof("error while reading request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading request"})
		return
	}
	jsonErr := json.Unmarshal(requestBytes, &req)
	if jsonErr != nil {
		logger.Log.Infof("error while marshalling json data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while marshalling json"})
		return
	}
	req.UserID = c.MustGet("userID").(uint)
	err = h.withdrawService.SaveWithdraw(req)
	if err != nil && err.Error() == "not enough funds" {
		logger.Log.Infof("not enough funds: %v", err)
		c.AbortWithStatus(http.StatusPaymentRequired)
		return
	} else if err != nil {
		logger.Log.Infof("error while saving withdraw: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving withdraw"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) GetAllWithdrawals(c RequestContext) {
	userID := c.MustGet("userID").(uint)
	withdrawals, err := h.withdrawService.GetAllWithdrawals(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if len(withdrawals) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.JSON(http.StatusOK, withdrawals)
}
