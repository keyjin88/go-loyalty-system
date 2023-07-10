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

func (r *OrderRepository) GetOrderByNumber(number string) (Order, error) {
	var saverOrder Order
	tx := r.db.First(&saverOrder, "number = ?", number)
	if tx.Error != nil {
		return Order{}, tx.Error
	}
	return saverOrder, nil
}
