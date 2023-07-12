package daemons

import (
	"encoding/json"
	"fmt"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/services"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

func WorkerProcessingOrders(ch <-chan storage.Order, host string, db *gorm.DB) {
	for order := range ch {
		logger.Log.Infof("processing %v", order)
		go func(order storage.Order) {
			getOrderDetails(&order, host)
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(order).Updates(order).Error; err != nil {
					return err
				}
				var savedUser storage.User
				if err := tx.First(&savedUser, "id = ?", order.UserID).Error; err != nil {
					return err
				}
				savedUser.Balance += order.Accrual
				if err := tx.Updates(&savedUser).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				logger.Log.Infof("Failed to process order %v: %v", order.ID, err)
			}
		}(order)
	}
}

func getOrderDetails(order *storage.Order, host string) {
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
