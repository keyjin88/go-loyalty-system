package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"io"
	"net/http"
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
	accrualServiceAddress string,
) *OrderService {
	go WorkerProcessingOrders(channel, accrualServiceAddress, orderRepository)
	return &OrderService{
		orderRepository:        orderRepository,
		orderProcessingChannel: channel,
	}
}

func (s *OrderService) SaveOrder(request storage.NewOrderRequest) (storage.Order, error) {
	if checkOrderNumber(request.Number) {
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
		} else {
			return storage.Order{}, err
		}
	}
	s.orderProcessingChannel <- order
	return order, nil
}

func WorkerProcessingOrders(ch <-chan storage.Order, host string, repository *storage.OrderRepository) {
	for order := range ch {
		logger.Log.Infof("processing %v", order)
		err := getOrderDetails(&order, host)
		if err != nil {
			logger.Log.Infof("error while processing: %e", err)
			return
		}
		err = repository.Save(&order)
		if err != nil {
			logger.Log.Infof("error while saving: %e", err)
			return
		}
	}
}

func getOrderDetails(order *storage.Order, host string) error {
	url := fmt.Sprintf(host+"/api/orders/%s", order.Number)
	maxRetries := 10
	retryInterval := 30 * time.Second
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		switch resp.StatusCode {
		case http.StatusOK:
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var details AccrualDetails
			err = json.Unmarshal(body, &details)
			if err != nil {
				return err
			}
			order.Status = details.Status
			order.Accrual = details.Accrual
		case http.StatusNoContent:
			return errors.New(fmt.Sprintf("Заказ %s не зарегистрирован в системе расчета", order.Number))
		case http.StatusTooManyRequests:
			if i == maxRetries-1 {
				return errors.New(fmt.Sprintf("Превышено количество запросов по заказу: %ы", order.Number))
			}
			resp.Body.Close()
			time.Sleep(retryInterval)
		case http.StatusInternalServerError:
			return errors.New(fmt.Sprintf("Внутренняя ошибка сервера"))
		default:
			return errors.New(fmt.Sprintf("Непредвиденный статус ответа: %s", resp.Status))
		}
	}
	return nil
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
