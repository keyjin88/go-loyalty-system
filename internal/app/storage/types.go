package storage

import (
	"gorm.io/gorm"
	"time"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type NewOrderRequest struct {
	Number string
	UserID uint
}

type Entity struct {
	gorm.Model
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
}

type User struct {
	Entity
	UserName  string  `json:"user_name" db:"user_name" gorm:"unique;not null"`
	Password  string  `json:"password" db:"password" gorm:"not null"`
	Balance   float64 `json:"balance" db:"balance" gorm:"default:0.0;not null"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn" gorm:"default:0.0;not null"`
}

type Order struct {
	Entity
	Number  string  `json:"number" db:"number" gorm:"unique;not null"`
	UserID  uint    `json:"user_id" db:"user_id" gorm:"not null"`
	Status  string  `json:"status" db:"status" gorm:"default:NEW;not null"`
	Accrual float64 `json:"accrual" db:"accrual"`
}

type AllOrderResponse struct {
	Number       string    `json:"number"`
	Status       string    `json:"status"`
	Accrual      float64   `json:"accrual"`
	UploadedDate time.Time `json:"-"`
	UploadedAt   string    `json:"uploaded_at"`
}

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type Withdraw struct {
	Entity
	OrderNumber string  `json:"order_number" db:"order_number" gorm:"not null"`
	Sum         float64 `json:"sum" db:"sum" gorm:"not null"`
	UserId      string  `json:"user_id" db:"user_id" gorm:"not null"`
}
