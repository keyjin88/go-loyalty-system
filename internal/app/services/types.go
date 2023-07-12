package services

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
)

//go:generate mockgen -destination=mocks/order_repository.go -package=mocks . OrderRepository
type OrderRepository interface {
	Save(order *storage.Order) error
	GetOrderByNumber(number string) (storage.Order, error)
	GetAllOrders(userID uint) ([]storage.Order, error)
}

//go:generate mockgen -destination=mocks/user_repository.go -package=mocks . UserRepository
type UserRepository interface {
	Save(user *storage.User) error
	Update(user *storage.User) error
	FindUserByID(userID uint) (storage.User, error)
	FindUserByUserName(userName string) (storage.User, error)
}

//go:generate mockgen -destination=mocks/withdraw_repository.go -package=mocks . WithdrawRepository
type WithdrawRepository interface {
	Save(withdraw *storage.Withdraw) error
	GetWithdrawals(userID uint) ([]storage.Withdraw, error)
}
