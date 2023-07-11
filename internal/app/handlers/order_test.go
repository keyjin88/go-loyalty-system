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
		saveOrderParameters storage.NewOrderRequest
		saveOrderResponse   storage.Order
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
			saveOrderParameters: storage.NewOrderRequest{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse: storage.Order{
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
			saveOrderParameters: storage.NewOrderRequest{},
			saveOrderResponse:   storage.Order{},
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
			saveOrderParameters: storage.NewOrderRequest{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  storage.Order{},
			saveOrderError:     errors.New("order already uploaded by this user"),
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
			saveOrderParameters: storage.NewOrderRequest{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  storage.Order{},
			saveOrderError:     errors.New("order already uploaded by another user"),
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
			saveOrderParameters: storage.NewOrderRequest{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  storage.Order{},
			saveOrderError:     errors.New("order has wrong format"),
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
			saveOrderParameters: storage.NewOrderRequest{
				Number: "1234567890",
				UserID: 101,
			},
			saveOrderResponse:  storage.Order{},
			saveOrderError:     errors.New("internal Server Error"),
			saveOrderCallCount: 1,
			status:             http.StatusInternalServerError,
			response:           gin.H{"error": "Internal Server Error"},
		},
	}
	for _, tt := range tests {
		orderService := mocks.NewMockOrderService(ctrl)
		requestContext := mocks.NewMockRequestContext(ctrl)

		requestContext.EXPECT().GetRawData().
			Return(tt.getRowData, tt.getRowDataError)
		requestContext.EXPECT().MustGet(tt.mustGetReturn).
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
