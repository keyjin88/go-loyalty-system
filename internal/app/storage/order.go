package storage

import (
	"gorm.io/gorm"
	"log"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	err := db.AutoMigrate(&Order{})
	if err != nil {
		log.Fatal("failed to migrate orders table")
	}
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Save(order *Order) error {
	err := r.db.Create(&order).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrderByNumber(number string) (Order, error) {
	var saverOrder Order
	tx := r.db.First(&saverOrder, "number = ?", number)
	if tx.Error != nil {
		return Order{}, tx.Error
	}
	return saverOrder, nil
}

func (r *OrderRepository) GetAllOrders(userID uint) ([]Order, error) {
	var orders []Order
	result := r.db.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
