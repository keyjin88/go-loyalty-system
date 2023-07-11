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
	MustGet(key string) any
}

//go:generate mockgen -destination=mocks/user_service.go -package=mocks . UserService
type UserService interface {
	SaveUser(request storage.AuthRequest) (storage.User, error)
	GetUserByUserName(request storage.AuthRequest) (storage.User, error)
	GetUserBalance(userID uint) (storage.BalanceResponse, error)
}

//go:generate mockgen -destination=mocks/order_service.go -package=mocks . OrderService
type OrderService interface {
	SaveOrder(orderNumber storage.NewOrderRequest) (storage.Order, error)
	GetAllOrders(userID uint) ([]storage.AllOrderResponse, error)
}

//go:generate mockgen -destination=mocks/withdraw_service.go -package=mocks . WithdrawService
type WithdrawService interface {
	SaveWithdraw(request storage.WithdrawRequest) error
	GetAllWithdrawals(userID uint) ([]storage.WithdrawResponse, error)
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
