package services

import (
	"errors"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"sort"
	"time"
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
		UserID:      user.ID,
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

func (s *WithdrawService) GetAllWithdrawals(userID uint) ([]storage.WithdrawResponse, error) {
	withdrawals, err := s.withdrawRepository.GetWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	var response = make([]storage.WithdrawResponse, 0)
	for _, withdraw := range withdrawals {
		resp := storage.WithdrawResponse{
			Order:         withdraw.OrderNumber,
			Sum:           withdraw.Sum,
			ProcessedDate: withdraw.CreatedAt,
			ProcessedAt:   withdraw.CreatedAt.Format(time.RFC3339),
		}
		response = append(response, resp)
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].ProcessedDate.Before(response[j].ProcessedDate)
	})
	return response, nil
}
