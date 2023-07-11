package storage

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"gorm.io/gorm"
	"log"
)

type OrderRepository struct {
	db                   *gorm.DB
	orderUpdatingChannel chan Order
}

func NewOrderRepository(db *gorm.DB, channel chan Order) *OrderRepository {
	err := db.AutoMigrate(&Order{})
	if err != nil {
		log.Fatal("failed to migrate orders table")
	}
	go WorkerUpdatingOrders(channel, db)
	return &OrderRepository{
		db:                   db,
		orderUpdatingChannel: channel,
	}
}

func (r *OrderRepository) Save(order *Order) error {
	err := r.db.Create(&order).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) Update(order *Order) {
	r.orderUpdatingChannel <- *order
}

func (r *OrderRepository) GetOrderByNumber(number string) (Order, error) {
	var saverOrder Order
	tx := r.db.First(&saverOrder, "number = ?", number)
	if tx.Error != nil {
		return Order{}, tx.Error
	}
	return saverOrder, nil
}

func (r *OrderRepository) GetAllOrders(userID int) ([]Order, error) {
	var orders []Order
	result := r.db.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func WorkerUpdatingOrders(ch <-chan Order, db *gorm.DB) {
	for order := range ch {
		err := db.Create(&order).Error
		if err != nil {
			logger.Log.Infof("error while updating order: %e", err)
			return
		}
		return
	}
}
