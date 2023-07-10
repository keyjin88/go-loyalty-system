package handlers

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
)

//go:generate mockgen -destination=mocks/request_context.go -package=mocks . RequestContext
type RequestContext interface {
	GetRawData() ([]byte, error)
	JSON(code int, obj any)
	AbortWithStatus(code int)
	Status(code int)
}

//go:generate mockgen -destination=mocks/user_service.go -package=mocks . UserService
type UserService interface {
	SaveUser(request storage.AuthRequest) (storage.User, error)
	GetUserByUserName(request storage.AuthRequest) (storage.User, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{
		userService: service,
	}
}
