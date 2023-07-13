package services

import (
	"errors"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/dto"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/models"
	"sort"
	"sync"
	"time"
)

type WithdrawService struct {
	withdrawRepository WithdrawRepository
	userRepository     UserRepository
	mutex              *sync.Mutex
}

func NewWithdrawService(
	withdrawRepository WithdrawRepository,
	userRepository UserRepository,
	mutex *sync.Mutex,
) *WithdrawService {
	return &WithdrawService{
		withdrawRepository: withdrawRepository,
		userRepository:     userRepository,
		mutex:              mutex,
	}
}

func (s *WithdrawService) SaveWithdraw(withdrawDTO dto.WithdrawDTO) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	user, err := s.userRepository.FindUserByID(withdrawDTO.UserID)
	if err != nil {
		return err
	}
	err = s.withdrawRepository.Save(&entities.Withdraw{
		OrderNumber: withdrawDTO.OrderNumber,
		Sum:         withdrawDTO.Sum,
		UserID:      user.ID,
	})
	if err != nil {
		return err
	}
	if user.Balance < withdrawDTO.Sum {
		return errors.New("not enough funds")
	}
	user.Balance -= withdrawDTO.Sum
	user.Withdrawn += withdrawDTO.Sum
	err = s.userRepository.Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (s *WithdrawService) GetAllWithdrawals(userID uint) ([]models.WithdrawResponse, error) {
	withdrawals, err := s.withdrawRepository.GetWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	var response = make([]models.WithdrawResponse, 0)
	for _, withdraw := range withdrawals {
		resp := models.WithdrawResponse{
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
