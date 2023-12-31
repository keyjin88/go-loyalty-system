package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers/mocks"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/dto"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/models"
	"net/http"
	"testing"
	"time"
)

func TestHandler_ProcessUserOrder(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                string
		getRowData          []byte
		getRowDataError     error
		mustGetReturn       uint
		mustGetCallCount    int
		saveOrderParameters dto.OrderDTO
		saveOrderResponse   entities.Order
		saveOrderError      error
		saveOrderCallCount  int
		status              int
		response            gin.H
	}{
		{
			name:             "Success",
			getRowData:       []byte("1234567890"),
			getRowDataError:  nil,
			mustGetReturn:    101,
			mustGetCallCount: 1,
			saveOrderParameters: dto.OrderDTO{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse: entities.Order{
				Number: "1234567890",
			},
			saveOrderError:     nil,
			saveOrderCallCount: 1,
			status:             http.StatusAccepted,
			response:           gin.H{"processed": "1234567890"},
		},
		{
			name:                "Error while reading request",
			getRowData:          nil,
			getRowDataError:     errors.New("error while reading request"),
			mustGetCallCount:    0,
			saveOrderParameters: dto.OrderDTO{},
			saveOrderResponse:   entities.Order{},
			saveOrderError:      nil,
			saveOrderCallCount:  0,
			status:              http.StatusBadRequest,
			response:            gin.H{"error": "Error while reading request"},
		},
		{
			name:             "Order already uploaded by this user",
			getRowData:       []byte("1234567890"),
			getRowDataError:  nil,
			mustGetReturn:    101,
			mustGetCallCount: 1,
			saveOrderParameters: dto.OrderDTO{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  entities.Order{},
			saveOrderError:     ErrOrderAlreadyUploadedByUser,
			saveOrderCallCount: 1,
			status:             http.StatusOK,
			response:           gin.H{"error": "order already uploaded by this user"},
		},
		{
			name:             "Order already uploaded by another user",
			getRowData:       []byte("1234567890"),
			getRowDataError:  nil,
			mustGetReturn:    101,
			mustGetCallCount: 1,
			saveOrderParameters: dto.OrderDTO{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  entities.Order{},
			saveOrderError:     ErrOrderAlreadyUploaded,
			saveOrderCallCount: 1,
			status:             http.StatusConflict,
			response:           gin.H{"error": "order already uploaded by another user"},
		},
		{
			name:             "Order has wrong format",
			getRowData:       []byte("1234567890"),
			getRowDataError:  nil,
			mustGetReturn:    101,
			mustGetCallCount: 1,
			saveOrderParameters: dto.OrderDTO{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  entities.Order{},
			saveOrderError:     ErrOrderHasWrongFormat,
			saveOrderCallCount: 1,
			status:             http.StatusUnprocessableEntity,
			response:           gin.H{"error": "wrong order number format"},
		},
		{
			name:             "Internal Server Error",
			getRowData:       []byte("1234567890"),
			getRowDataError:  nil,
			mustGetReturn:    101,
			mustGetCallCount: 1,
			saveOrderParameters: dto.OrderDTO{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  entities.Order{},
			saveOrderError:     errors.New("internal Server Error"),
			saveOrderCallCount: 1,
			status:             http.StatusInternalServerError,
			response:           gin.H{"error": "Internal Server Error"},
		},
	}
	orderService := mocks.NewMockOrderService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		requestContext.EXPECT().GetRawData().
			Return(tt.getRowData, tt.getRowDataError)
		requestContext.EXPECT().MustGet("userID").
			Return(tt.mustGetReturn).
			Times(tt.mustGetCallCount)
		requestContext.EXPECT().JSON(tt.status, tt.response)

		orderService.EXPECT().SaveOrder(tt.saveOrderParameters).
			Return(tt.saveOrderResponse, tt.saveOrderError).
			Times(tt.saveOrderCallCount)

		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				orderService: orderService,
			}
			h.ProcessUserOrder(requestContext)
		})
	}
}

func TestHandler_GetAllOrders(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	orders := []models.AllOrderResponse{
		{
			Number:       "111111111",
			Status:       "NEW",
			Accrual:      123.32,
			UploadedDate: now,
			UploadedAt:   now.Format(time.RFC3339),
		},
	}

	tests := []struct {
		name               string
		mustGetReturn      uint
		mustGetCallCount   int
		getAllOrdersReturn []models.AllOrderResponse
		getAllOrdersError  error
		status             int
		response           any
	}{
		{
			name:               "Success",
			mustGetReturn:      101,
			mustGetCallCount:   1,
			getAllOrdersReturn: orders,
			getAllOrdersError:  nil,
			status:             http.StatusOK,
			response:           orders,
		},
		{
			name:               "Internal Server Error",
			mustGetReturn:      101,
			mustGetCallCount:   1,
			getAllOrdersReturn: []models.AllOrderResponse{},
			getAllOrdersError:  errors.New("internal Server Error"),
			status:             http.StatusInternalServerError,
			response:           gin.H{"error": "Internal Server Error"},
		},
		{
			name:               "Orders not found",
			mustGetReturn:      101,
			mustGetCallCount:   1,
			getAllOrdersReturn: []models.AllOrderResponse{},
			getAllOrdersError:  nil,
			status:             http.StatusNoContent,
			response:           gin.H{"error": "orders not found"},
		},
	}
	orderService := mocks.NewMockOrderService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext.EXPECT().MustGet(gomock.Any()).
				Return(tt.mustGetReturn).
				Times(tt.mustGetCallCount)
			requestContext.EXPECT().JSON(tt.status, tt.response)

			orderService.EXPECT().GetAllOrders(gomock.Any()).Return(tt.getAllOrdersReturn, tt.getAllOrdersError)

			h := &Handler{
				orderService: orderService,
			}
			h.GetAllOrders(requestContext)
		})
	}
}
