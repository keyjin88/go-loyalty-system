package handlers

import "net/http"

func (h *Handler) GetBalance(c RequestContext) {
	userID := c.MustGet("userID").(uint)
	response, err := h.userService.GetUserBalance(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, response)
}
