package dto

type OrderDTO struct {
	Number string
	UserID uint
}

type UserDTO struct {
	UserName string
	Password string
}

type WithdrawDTO struct {
	OrderNumber string
	Sum         float64
	UserID      uint
}
