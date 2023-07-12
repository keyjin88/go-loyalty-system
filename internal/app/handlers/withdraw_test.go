package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers/mocks"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"net/http"
	"testing"
	"time"
)

func TestHandler_SaveWithdraw(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                  string
		getRowDataReturn      []byte
		getRowDataError       error
		mustGetReturn         uint
		mustGetCallCount      int
		saveWithdrawError     error
		saveWithdrawCallCount int
		status                int
		response              gin.H
	}{
		{
			name:                  "Success",
			getRowDataReturn:      []byte("{\n\"order\": \"2377225626\",\n\"sum\": 1000\n}"),
			getRowDataError:       nil,
			mustGetReturn:         101,
			mustGetCallCount:      1,
			saveWithdrawError:     nil,
			saveWithdrawCallCount: 1,
			status:                http.StatusOK,
			response:              gin.H{"info": "Withdrawal successfully saved"},
		},
		{
			name:                  "Error while reading request",
			getRowDataReturn:      nil,
			getRowDataError:       errors.New("error while reading request"),
			mustGetReturn:         101,
			mustGetCallCount:      0,
			saveWithdrawError:     nil,
			saveWithdrawCallCount: 0,
			status:                http.StatusBadRequest,
			response:              gin.H{"error": "Error while reading request"},
		},
		{
			name:                  "Error while marshalling json",
			getRowDataReturn:      []byte("BAD JSON"),
			getRowDataError:       nil,
			mustGetReturn:         101,
			mustGetCallCount:      0,
			saveWithdrawError:     nil,
			saveWithdrawCallCount: 0,
			status:                http.StatusBadRequest,
			response:              gin.H{"error": "Error while marshalling json"},
		},
		{
			name:                  "Not enough funds",
			getRowDataReturn:      []byte("{\n\"order\": \"2377225626\",\n\"sum\": 1000\n}"),
			getRowDataError:       nil,
			mustGetReturn:         101,
			mustGetCallCount:      1,
			saveWithdrawError:     errors.New("not enough funds"),
			saveWithdrawCallCount: 1,
			status:                http.StatusPaymentRequired,
			response:              gin.H{"error": "Not enough funds"},
		},
		{
			name:                  "Error while saving withdraw",
			getRowDataReturn:      []byte("{\n\"order\": \"2377225626\",\n\"sum\": 1000\n}"),
			getRowDataError:       nil,
			mustGetReturn:         101,
			mustGetCallCount:      1,
			saveWithdrawError:     errors.New("error while saving withdraw"),
			saveWithdrawCallCount: 1,
			status:                http.StatusInternalServerError,
			response:              gin.H{"error": "Error while saving withdraw"},
		},
	}
	withdrawService := mocks.NewMockWithdrawService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		withdrawService.EXPECT().SaveWithdraw(gomock.Any()).
			Return(tt.saveWithdrawError).
			Times(tt.saveWithdrawCallCount)
		requestContext.EXPECT().MustGet("userID").
			Return(tt.mustGetReturn).
			Times(tt.mustGetCallCount)
		requestContext.EXPECT().GetRawData().Return(tt.getRowDataReturn, tt.getRowDataError)
		requestContext.EXPECT().JSON(tt.status, tt.response)

		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				withdrawService: withdrawService,
			}
			h.SaveWithdraw(requestContext)
		})
	}
}

func TestHandler_GetAllWithdrawals(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	withdrawals := []storage.WithdrawResponse{
		{
			Order:         "123",
			Sum:           123,
			ProcessedDate: now,
			ProcessedAt:   now.Format(time.RFC3339),
		},
	}

	tests := []struct {
		name                    string
		mustGetReturn           uint
		mustGetCallCount        int
		getAllWithdrawalsReturn []storage.WithdrawResponse
		getAllWithdrawalsError  error
		status                  int
		response                any
	}{
		{
			name:                    "Success",
			mustGetReturn:           101,
			mustGetCallCount:        1,
			getAllWithdrawalsReturn: withdrawals,
			getAllWithdrawalsError:  nil,
			status:                  http.StatusOK,
			response:                withdrawals,
		},
		{
			name:                    "Internal server error",
			mustGetReturn:           101,
			mustGetCallCount:        1,
			getAllWithdrawalsReturn: withdrawals,
			getAllWithdrawalsError:  errors.New("internal server error"),
			status:                  http.StatusInternalServerError,
			response:                gin.H{"error": "Internal server error"},
		},
		{
			name:                    "Withdrawals are empty",
			mustGetReturn:           101,
			mustGetCallCount:        1,
			getAllWithdrawalsReturn: []storage.WithdrawResponse{},
			getAllWithdrawalsError:  nil,
			status:                  http.StatusNoContent,
			response:                gin.H{"error": "withdrawal not found"},
		},
	}
	withdrawService := mocks.NewMockWithdrawService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		requestContext.EXPECT().MustGet(gomock.Any()).
			Return(tt.mustGetReturn).
			Times(tt.mustGetCallCount)
		requestContext.EXPECT().JSON(tt.status, tt.response)
		withdrawService.EXPECT().GetAllWithdrawals(tt.mustGetReturn).
			Return(tt.getAllWithdrawalsReturn, tt.getAllWithdrawalsError)
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				withdrawService: withdrawService,
			}
			h.GetAllWithdrawals(requestContext)
		})
	}
}
