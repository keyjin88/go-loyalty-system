package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers/mocks"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"testing"
)

func TestHandler_RegisterUser(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := mocks.NewMockUserService(ctrl)
			requestContext := mocks.NewMockRequestContext(ctrl)

			h := &Handler{
				userService: userService,
			}
			h.RegisterUser(requestContext)
		})
	}
}
