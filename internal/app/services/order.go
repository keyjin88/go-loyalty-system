package services

import (
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/dto"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AccrualDetails struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type OrderService struct {
	orderRepository        OrderRepository
	orderProcessingChannel chan entities.Order
}

func NewOrderService(
	orderRepository OrderRepository,
	channel chan entities.Order,
) *OrderService {
	return &OrderService{
		orderRepository:        orderRepository,
		orderProcessingChannel: channel,
	}
}

func (s *OrderService) SaveOrder(orderDTO dto.OrderDTO) (entities.Order, error) {
	// закомментировано, длч облегчния тестирования
	if !checkOrderNumber(orderDTO.Number) {
		return entities.Order{}, handlers.ErrOrderHasWrongFormat
	}
	var order = entities.Order{
		Number: orderDTO.Number,
		UserID: orderDTO.UserID,
	}
	err := s.orderRepository.Save(&order)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			order, err = s.orderRepository.GetOrderByNumber(orderDTO.Number)
			if err != nil {
				return entities.Order{}, err
			}
			if order.UserID == orderDTO.UserID {
				return entities.Order{}, handlers.ErrOrderAlreadyUploadedByUser
			}
			return entities.Order{}, handlers.ErrOrderAlreadyUploaded
		}
		return entities.Order{}, err
	}
	s.orderProcessingChannel <- order
	return order, nil
}

func (s *OrderService) GetAllOrders(userID uint) ([]models.AllOrderResponse, error) {
	orders, err := s.orderRepository.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}
	var response = make([]models.AllOrderResponse, 0)
	for _, order := range orders {
		resp := models.AllOrderResponse{
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
	// Удаляем все пробелы из строки
	orderNumber = strings.ReplaceAll(orderNumber, " ", "")

	// Проверяем, что номер заказа состоит только из цифр
	_, err := strconv.Atoi(orderNumber)
	if err != nil {
		return false
	}

	// Применяем алгоритм Луна для валидации номера заказа
	sum := 0
	double := false
	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(orderNumber[i]))

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	return sum%10 == 0
}
