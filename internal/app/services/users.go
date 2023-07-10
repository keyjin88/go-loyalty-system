package services

import (
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	userRepository *storage.UserRepository
}

func NewUserService(userRepository *storage.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) SaveUser(request storage.RegisterUserRequest) (storage.User, error) {
	//зашифровать пароль
	user := storage.User{
		UserName: request.Login,
		Password: hashPassword(request.Password),
	}
	// Сохранить юзера в БД
	err := s.userRepository.Save(&user)
	if err != nil {
		return storage.User{}, err
	}
	return user, nil
}

// Хэширование пароля
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash password")
	}
	return string(hash)
}
