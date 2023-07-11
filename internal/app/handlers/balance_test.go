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
)

func TestHandler_GetBalance(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name              string
		userID            uint
		userServiceReturn storage.BalanceResponse
		userServiceError  error
		status            int
		response          any
	}{
		{
			name:   "Success",
			userID: 101,
			userServiceReturn: storage.BalanceResponse{
				Current:   101.01,
				Withdrawn: 543.21,
			},
			userServiceError: nil,
			status:           http.StatusOK,
			response: storage.BalanceResponse{
				Current:   101.01,
				Withdrawn: 543.21,
			},
		},
		{
			name:              "Internal Server Error",
			userID:            101,
			userServiceReturn: storage.BalanceResponse{},
			userServiceError:  errors.New("error"),
			status:            http.StatusInternalServerError,
			response:          gin.H{"error": "Internal Server Error"},
		},
	}
	for _, tt := range tests {
		userService := mocks.NewMockUserService(ctrl)
		userService.EXPECT().GetUserBalance(tt.userID).Return(tt.userServiceReturn, tt.userServiceError)
		requestContext := mocks.NewMockRequestContext(ctrl)
		requestContext.EXPECT().MustGet("mustGetReturn").Return(tt.userID)
		requestContext.EXPECT().JSON(tt.status, tt.response)

		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				userService: userService,
			}
			h.GetBalance(requestContext)
		})
	}
}
