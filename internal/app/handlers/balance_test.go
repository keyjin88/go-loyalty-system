package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers/mocks"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/models"
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
		mustGetReturn     uint
		userServiceReturn models.BalanceResponse
		userServiceError  error
		status            int
		response          any
	}{
		{
			name:          "Success",
			mustGetReturn: 101,
			userServiceReturn: models.BalanceResponse{
				Current:   101.01,
				Withdrawn: 543.21,
			},
			userServiceError: nil,
			status:           http.StatusOK,
			response: models.BalanceResponse{
				Current:   101.01,
				Withdrawn: 543.21,
			},
		},
		{
			name:              "Internal Server Error",
			mustGetReturn:     101,
			userServiceReturn: models.BalanceResponse{},
			userServiceError:  errors.New("error"),
			status:            http.StatusInternalServerError,
			response:          gin.H{"error": "Internal Server Error"},
		},
	}
	userService := mocks.NewMockUserService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		userService.EXPECT().GetUserBalance(tt.mustGetReturn).
			Return(tt.userServiceReturn, tt.userServiceError)
		requestContext.EXPECT().MustGet(gomock.Any()).
			Return(tt.mustGetReturn)
		requestContext.EXPECT().JSON(tt.status, tt.response)
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				userService: userService,
			}
			h.GetBalance(requestContext)
		})
	}
}
