package entities

import "gorm.io/gorm"

type Entity struct {
	gorm.Model
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
}

type Withdraw struct {
	Entity
	OrderNumber string  `json:"order_number" db:"order_number" gorm:"not null"`
	Sum         float64 `json:"sum" db:"sum" gorm:"not null"`
	UserID      uint    `json:"user_id" db:"user_id" gorm:"not null"`
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
