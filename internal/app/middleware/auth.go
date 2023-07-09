package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	// Получаем токен из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не предоставлен токен авторизации"})
		c.Abort()
		return
	}

	// Проверяем формат токена
	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что алгоритм подписи совпадает с ожидаемым
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неверный алгоритм подписи: %v", token.Header["alg"])
		}

		// Возвращаем ключ для проверки подписи
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен авторизации"})
		c.Abort()
		return
	}

	// Проверяем, что токен является объектом типа JWT и содержит правильные поля
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Получаем userID из токена
		userID := claims["userID"].(string)
		c.Set("userID", userID)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен авторизации"})
		c.Abort()
		return
	}
}