package handlers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"net/http"
	"time"
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
	userFromDB, err := h.userService.SaveUser(req)
	if err != nil {
		if err.Error() == "user already exists" {
			logger.Log.Infof("User already exists: %v", err)
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		} else {
			logger.Log.Infof("Internal Server Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
	token, err := createToken(userFromDB.ID, h.secret)
	if err != nil {
		logger.Log.Infof("Failed to create token %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create JWT token"})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"info": "New user successfully created"})
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
	userFromDB, err := h.userService.GetUserByUserName(req)
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			logger.Log.Infof("Wrong password %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}
		logger.Log.Infof("Failed to get user %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user by username"})
		return
	}
	token, err := createToken(userFromDB.ID, h.secret)
	if err != nil {
		logger.Log.Infof("Failed to create token %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create JWT token"})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"info": "login successful"})
}

func createToken(userID uint, secret string) (string, error) {
	if userID == 0 || len(secret) == 0 {
		return "", errors.New("invalid token credentials")
	}
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Токен действителен в течение 24 часов
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
