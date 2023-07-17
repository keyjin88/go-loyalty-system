package models

import "time"

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type NewOrderRequest struct {
	Number string
	UserID uint
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
	UserID uint    `json:"-"`
	Order  string  `json:"order"`
	Sum    float64 `json:"sum"`
}

type WithdrawResponse struct {
	Order         string    `json:"order"`
	Sum           float64   `json:"sum"`
	ProcessedDate time.Time `json:"-"`
	ProcessedAt   string    `json:"processed_at"`
}
