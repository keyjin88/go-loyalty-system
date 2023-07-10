package storage

import "time"

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type NewOrderRequest struct {
	Number string
	UserID int
}

type Entity struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}

type User struct {
	Entity
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
}

type Order struct {
	Entity
	Number string `json:"number" db:"number"`
	UserID int    `json:"user_id" db:"user_id"`
}
