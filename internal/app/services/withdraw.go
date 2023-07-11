package services

import (
	"errors"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
)

type WithdrawService struct {
	withdrawRepository *storage.WithdrawRepository
	userRepository     *storage.UserRepository
}

func NewWithdrawService(withdrawRepository *storage.WithdrawRepository, userRepository *storage.UserRepository) *WithdrawService {
	return &WithdrawService{
		withdrawRepository: withdrawRepository,
		userRepository:     userRepository,
	}
}

func (s *WithdrawService) SaveWithdraw(request storage.WithdrawRequest) error {
	user, err := s.userRepository.FindUserByID(request.UserID)
	if err != nil {
		return err
	}
	err = s.withdrawRepository.Save(&storage.Withdraw{
		OrderNumber: request.Order,
		Sum:         request.Sum,
		UserId:      user.ID,
	})
	if err != nil {
		return err
	}
	if user.Balance < request.Sum {
		return errors.New("not enough funds")
	}
	user.Balance -= request.Sum
	user.Withdrawn += request.Sum
	err = s.userRepository.Update(&user)
	if err != nil {
		return err
	}
	return nil
}
