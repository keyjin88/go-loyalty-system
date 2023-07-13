package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/dto"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/models"
)

//go:generate mockgen -destination=mocks/request_context.go -package=mocks . RequestContext
type RequestContext interface {
	GetRawData() ([]byte, error)
	JSON(code int, obj any)
	Header(key, value string)
	MustGet(key string) any
}

//go:generate mockgen -destination=mocks/user_service.go -package=mocks . UserService
type UserService interface {
	SaveUser(userDTO dto.UserDTO) (entities.User, error)
	GetUserByUserName(userDTO dto.UserDTO) (entities.User, error)
	GetUserBalance(userID uint) (models.BalanceResponse, error)
}

//go:generate mockgen -destination=mocks/order_service.go -package=mocks . OrderService
type OrderService interface {
	SaveOrder(orderNumber dto.OrderDTO) (entities.Order, error)
	GetAllOrders(userID uint) ([]models.AllOrderResponse, error)
}

//go:generate mockgen -destination=mocks/withdraw_service.go -package=mocks . WithdrawService
type WithdrawService interface {
	SaveWithdraw(withdrawDTO dto.WithdrawDTO) error
	GetAllWithdrawals(userID uint) ([]models.WithdrawResponse, error)
}

type Claims struct {
	UserID uint `json:"userID"`
	jwt.StandardClaims
}

type Handler struct {
	userService     UserService
	orderService    OrderService
	withdrawService WithdrawService
	secret          string
}

func NewHandler(userService UserService, oderService OrderService, withdrawService WithdrawService, secret string) *Handler {
	return &Handler{
		userService:     userService,
		orderService:    oderService,
		withdrawService: withdrawService,
		secret:          secret,
	}
}
