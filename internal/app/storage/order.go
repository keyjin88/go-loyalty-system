package storage

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"gorm.io/gorm"
	"log"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	err := db.AutoMigrate(&entities.Order{})
	if err != nil {
		log.Fatal("failed to migrate orders table")
	}
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Save(order *entities.Order) error {
	err := r.db.Create(&order).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrderByNumber(number string) (entities.Order, error) {
	var saverOrder entities.Order
	tx := r.db.First(&saverOrder, "number = ?", number)
	if tx.Error != nil {
		return entities.Order{}, tx.Error
	}
	return saverOrder, nil
}

func (r *OrderRepository) GetAllOrders(userID uint) ([]entities.Order, error) {
	var orders []entities.Order
	result := r.db.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
