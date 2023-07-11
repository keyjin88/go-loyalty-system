package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"log"
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
			c.AbortWithStatus(http.StatusConflict)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	token, err := createToken(userFromDB.ID, h.secret)
	if err != nil {
		log.Fatal("failed to create JWT token")
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, "New user successfully created")
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
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return
	}
	token, err := createToken(userFromDB.ID, h.secret)
	if err != nil {
		log.Fatal("failed to create JWT token")
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, "Login successful")
}

func createToken(userID uint, secret string) (string, error) {
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
