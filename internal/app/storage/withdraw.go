package storage

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"gorm.io/gorm"
	"log"
)

type WithdrawRepository struct {
	db *gorm.DB
}

func NewWithdrawRepository(db *gorm.DB) *WithdrawRepository {
	err := db.AutoMigrate(&entities.Withdraw{})
	if err != nil {
		log.Fatal("failed to migrate withdraw table")
	}
	return &WithdrawRepository{
		db: db,
	}
}

func (r *WithdrawRepository) Save(withdraw *entities.Withdraw) error {
	err := r.db.Create(&withdraw).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WithdrawRepository) GetWithdrawals(userID uint) ([]entities.Withdraw, error) {
	var withdraws []entities.Withdraw
	result := r.db.Where("user_id = ?", userID).Find(&withdraws)
	if result.Error != nil {
		return nil, result.Error
	}
	return withdraws, nil
}
