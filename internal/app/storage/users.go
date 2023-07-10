package storage

import (
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(user *User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		log.Fatal("failed to save user to database")
	}
	return nil
}
