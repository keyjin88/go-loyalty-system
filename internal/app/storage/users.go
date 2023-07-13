package storage

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	err := db.AutoMigrate(&entities.User{})
	if err != nil {
		log.Fatal("failed to migrate users table")
	}
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(user *entities.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(user *entities.User) error {
	err := r.db.Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByID(userID uint) (entities.User, error) {
	var savedUser entities.User
	tx := r.db.First(&savedUser, "id = ?", userID)
	if tx.Error != nil {
		return entities.User{}, tx.Error
	}
	return savedUser, nil
}

func (r *UserRepository) FindUserByUserName(userName string) (entities.User, error) {
	var savedUser entities.User
	tx := r.db.First(&savedUser, "user_name = ?", userName)
	if tx.Error != nil {
		return entities.User{}, tx.Error
	}
	return savedUser, nil
}
