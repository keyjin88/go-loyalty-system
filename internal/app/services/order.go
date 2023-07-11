package services

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"sort"
	"strconv"
	"time"
)

type AccrualDetails struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type OrderService struct {
	orderRepository        *storage.OrderRepository
	orderProcessingChannel chan storage.Order
}

func NewOrderService(
	orderRepository *storage.OrderRepository,
	channel chan storage.Order,
) *OrderService {
	return &OrderService{
		orderRepository:        orderRepository,
		orderProcessingChannel: channel,
	}
}

func (s *OrderService) SaveOrder(request storage.NewOrderRequest) (storage.Order, error) {
	if !checkOrderNumber(request.Number) {
		return storage.Order{}, errors.New("order has wrong format")
	}
	var order = storage.Order{
		Number: request.Number,
		UserID: request.UserID,
	}
	err := s.orderRepository.Save(&order)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			order, err = s.orderRepository.GetOrderByNumber(request.Number)
			if err != nil {
				return storage.Order{}, err
			}
			if order.UserID == request.UserID {
				return storage.Order{}, errors.New("order already uploaded by this user")
			}
			return storage.Order{}, errors.New("order already uploaded by another user")
		}
		return storage.Order{}, err
	}
	s.orderProcessingChannel <- order
	return order, nil
}

func (s *OrderService) GetAllOrders(userID int) ([]storage.AllOrderResponse, error) {
	orders, err := s.orderRepository.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}
	var response = make([]storage.AllOrderResponse, 0)
	for _, order := range orders {
		resp := storage.AllOrderResponse{
			Number:       order.Number,
			Status:       order.Status,
			Accrual:      order.Accrual,
			UploadedDate: order.CreatedAt,
			UploadedAt:   order.CreatedAt.Format(time.RFC3339),
		}
		response = append(response, resp)
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].UploadedDate.Before(response[j].UploadedDate)
	})
	return response, nil
}

func checkOrderNumber(orderNumber string) bool {
	// Удаляем все пробелы из номера заказа
	orderNumber = removeSpaces(orderNumber)
	// Проверяем, что номер заказа состоит только из цифр
	_, err := strconv.Atoi(orderNumber)
	if err != nil {
		return false
	}
	// Проверяем длину номера заказа
	if len(orderNumber) < 9 || len(orderNumber) > 16 {
		return false
	}
	// Вычисляем контрольную сумму по алгоритму Луна
	sum := 0
	for i, digit := range orderNumber {
		// Преобразуем символ цифры в число
		num, _ := strconv.Atoi(string(digit))
		// Удваиваем каждую вторую цифру, начиная с последней
		if i%2 == len(orderNumber)%2 {
			num *= 2
			// Если результат удвоения больше 9, вычитаем 9
			if num > 9 {
				num -= 9
			}
		}
		// Суммируем все цифры
		sum += num
	}
	// Проверяем, что контрольная сумма делится нацело на 10
	return sum%10 == 0
}

func removeSpaces(s string) string {
	result := ""
	for _, char := range s {
		if char != ' ' {
			result += string(char)
		}
	}
	return result
}
