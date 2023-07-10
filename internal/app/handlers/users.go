package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"net/http"
)

func (h *Handler) RegisterUser(c RequestContext) {
	var req storage.AuthRequest
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
	_, err = h.userService.SaveUser(req)
	if err != nil {
		if err.Error() == "user already exists" {
			c.AbortWithStatus(http.StatusConflict)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) LoginUser(c RequestContext) {
	var req storage.AuthRequest
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
	_, err = h.userService.GetUserByUserName(req)
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return
	}
	c.Status(http.StatusCreated)
}
