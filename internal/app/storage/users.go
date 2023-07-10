package storage

import (
	"gorm.io/gorm"
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
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByUserName(userName string) (User, error) {
	var savedUser User
	tx := r.db.First(&savedUser, "user_name = ?", userName)
	if tx.Error != nil {
		return User{}, tx.Error
	}
	return savedUser, nil
}
