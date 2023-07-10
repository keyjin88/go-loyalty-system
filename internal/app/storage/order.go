package storage

import "gorm.io/gorm"

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
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
