package services

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *UserService) SaveUser(request storage.AuthRequest) (storage.User, error) {
	user := storage.User{
		UserName: request.Login,
		Password: hashPassword(request.Password),
	}
	err := s.userRepository.Save(&user)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			return storage.User{}, errors.New("user already exists")
		} else {
			return storage.User{}, err
		}
	}
	return user, nil
}

func (s *UserService) GetUserByUserName(request storage.AuthRequest) (storage.User, error) {
	user, err := s.userRepository.FindUserByUserName(request.Login)
	if err != nil {
		return storage.User{}, err
	}
	passwordError := comparePassword(user.Password, request.Password)
	if passwordError != nil {
		return storage.User{}, passwordError
	}
	return user, nil
}

func (s *UserService) GetUserBalance(userID uint) (storage.BalanceResponse, error) {
	user, err := s.userRepository.FindUserById(userID)
	if err != nil {
		return storage.BalanceResponse{}, err
	}
	return storage.BalanceResponse{Current: user.Balance, Withdrawn: user.Withdrawn}, nil
}

// Хэширование пароля
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash password")
	}
	return string(hash)
}

func comparePassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
