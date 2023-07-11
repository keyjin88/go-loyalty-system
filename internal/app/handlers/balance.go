package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetBalance(c RequestContext) {
	userID := c.MustGet("userID").(uint)
	response, err := h.userService.GetUserBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	c.JSON(http.StatusOK, response)
}
