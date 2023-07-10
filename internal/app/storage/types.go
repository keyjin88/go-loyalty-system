package storage

import "time"

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserName  string    `json:"user_name" db:"user_name"`
	Password  string    `json:"password" db:"password"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
