package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
)

//go:generate mockgen -destination=mocks/request_context.go -package=mocks . RequestContext
type RequestContext interface {
	GetRawData() ([]byte, error)
	JSON(code int, obj any)
	AbortWithStatus(code int)
	Status(code int)
	Header(key, value string)
}

//go:generate mockgen -destination=mocks/user_service.go -package=mocks . UserService
type UserService interface {
	SaveUser(request storage.AuthRequest) (storage.User, error)
	GetUserByUserName(request storage.AuthRequest) (storage.User, error)
}

type Claims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}

type Handler struct {
	userService UserService
	secret      string
}

func NewHandler(service UserService, secret string) *Handler {
	return &Handler{
		userService: service,
		secret:      secret,
	}
}
