package handlers

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/services"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
)

//go:generate mockgen -destination=mocks/request_context.go -package=mocks . RequestContext
type RequestContext interface {
	GetRawData() ([]byte, error)
	JSON(code int, obj any)
}

//go:generate mockgen -destination=mocks/user_service.go -package=mocks . UserService
type UserService interface {
	SaveUser(request storage.RegisterUserRequest) (storage.User, error)
}

type Handler struct {
	userService *services.UserService
}

func NewHandler(userService *services.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}
