package daemons

import (
	"encoding/json"
	"fmt"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/services"
	"gorm.io/gorm"
	"io"
	"net/http"
	"sync"
	"time"
)

func WorkerProcessingOrders(ch <-chan entities.Order, host string, db *gorm.DB, maxWorkers int, mutex *sync.Mutex) {
	workerPool := make(chan struct{}, maxWorkers) // Создаем пул горутин
	for order := range ch {
		workerPool <- struct{}{} // Заполняем пул горутин
		go func(orderID uint, accrual float64) {
			mutex.Lock()
			defer mutex.Unlock()
			defer func() {
				<-workerPool // Освобождаем горутину при завершении
			}()
			var order entities.Order
			if err := db.First(&order, orderID).Error; err != nil {
				logger.Log.Errorf("Failed to retrieve order %v: %v", orderID, err)
				return
			}
			getOrderDetails(&order, host)
			var savedUser entities.User
			if err := db.First(&savedUser, "id = ?", order.UserID).Error; err != nil {
				logger.Log.Errorf("Failed to retrieve user %v: %v", order.UserID, err)
				return
			}
			err := db.Set("gorm:query_option", "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE").
				Transaction(func(tx *gorm.DB) error {
					if err := tx.Model(order).Updates(order).Error; err != nil {
						return err
					}
					var savedUser entities.User
					if err := tx.First(&savedUser, "id = ?", order.UserID).Error; err != nil {
						return err
					}
					savedUser.Balance += order.Accrual
					return tx.Updates(&savedUser).Error
				})
			if err != nil {
				logger.Log.Error("Failed to process order %v: %v", order.ID, err)
			}
		}(order.ID, order.Accrual)
	}
}

func getOrderDetails(order *entities.Order, host string) {
	url := fmt.Sprintf(host+"/api/orders/%s", order.Number)
	maxRetries := 5
	retryInterval := 1 * time.Second
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Infof("Error getting order info from: %s", url)
			return
		}
		switch resp.StatusCode {
		case http.StatusOK:
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Log.Infof("Error reading response")
				return
			}
			var details services.AccrualDetails
			err = json.Unmarshal(body, &details)
			if err != nil {
				logger.Log.Infof("Error unmarshalling response")
				return
			}
			order.Status = details.Status
			order.Accrual = details.Accrual
		case http.StatusNoContent:
			logger.Log.Infof("заказ %s не зарегистрирован в системе расчета", order.Number)
			return
		case http.StatusTooManyRequests:
			if i == maxRetries-1 {
				logger.Log.Infof("превышено количество запросов по заказу: %s", order.Number)
				return
			}
			err := resp.Body.Close()
			if err != nil {
				logger.Log.Infof("error while closing response body")
				return
			}
			time.Sleep(retryInterval)
		case http.StatusInternalServerError:
			logger.Log.Infof("внутренняя ошибка сервера")
			return
		default:
			logger.Log.Infof("непредвиденный статус ответа: %s", resp.Status)
			return
		}
	}
}
