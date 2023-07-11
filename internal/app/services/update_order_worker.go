package services

import (
	"encoding/json"
	"fmt"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

func WorkerProcessingOrders(ch <-chan storage.Order, host string, db *gorm.DB) {
	for order := range ch {
		logger.Log.Infof("processing %v", order)
		go func(order storage.Order) {
			getOrderDetails(&order, host)
			err := db.Model(&order).Updates(order).Error
			if err != nil {
				log.Printf("Failed to update order %v: %v", order.ID, err)
			} else {
				log.Printf("Order %v updated", order.ID)
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
			var details AccrualDetails
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