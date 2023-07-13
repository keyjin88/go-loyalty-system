package services

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/dto"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) SaveUser(userDTO dto.UserDTO) (entities.User, error) {
	user := entities.User{
		UserName: userDTO.UserName,
		Password: hashPassword(userDTO.Password),
	}
	err := s.userRepository.Save(&user)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			return entities.User{}, errors.New("user already exists")
		} else {
			return entities.User{}, err
		}
	}
	return user, nil
}

func (s *UserService) GetUserByUserName(userDTO dto.UserDTO) (entities.User, error) {
	user, err := s.userRepository.FindUserByUserName(userDTO.UserName)
	if err != nil {
		return entities.User{}, err
	}
	passwordError := comparePassword(user.Password, userDTO.Password)
	if passwordError != nil {
		return entities.User{}, passwordError
	}
	return user, nil
}

func (s *UserService) GetUserBalance(userID uint) (models.BalanceResponse, error) {
	user, err := s.userRepository.FindUserByID(userID)
	if err != nil {
		return models.BalanceResponse{}, err
	}
	return models.BalanceResponse{Current: user.Balance, Withdrawn: user.Withdrawn}, nil
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
