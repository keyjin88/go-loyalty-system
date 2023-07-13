package services

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
)

//go:generate mockgen -destination=mocks/order_repository.go -package=mocks . OrderRepository
type OrderRepository interface {
	Save(order *entities.Order) error
	GetOrderByNumber(number string) (entities.Order, error)
	GetAllOrders(userID uint) ([]entities.Order, error)
}

//go:generate mockgen -destination=mocks/user_repository.go -package=mocks . UserRepository
type UserRepository interface {
	Save(user *entities.User) error
	Update(user *entities.User) error
	FindUserByID(userID uint) (entities.User, error)
	FindUserByUserName(userName string) (entities.User, error)
}

//go:generate mockgen -destination=mocks/withdraw_repository.go -package=mocks . WithdrawRepository
type WithdrawRepository interface {
	Save(withdraw *entities.Withdraw) error
	GetWithdrawals(userID uint) ([]entities.Withdraw, error)
}
