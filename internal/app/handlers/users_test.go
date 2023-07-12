package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers/mocks"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"gorm.io/gorm"
	"net/http"
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
		name              string
		getRowDataReturn  []byte
		getRowDataError   error
		saveUserReturn    storage.User
		saveUserError     error
		saveUserCallCount int
		headerCallCount   int
		status            int
		response          gin.H
		secret            string
	}{
		{
			name:             "Success",
			getRowDataReturn: []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:  nil,
			saveUserReturn: storage.User{
				Entity: storage.Entity{
					Model:     gorm.Model{ID: 1},
					IsDeleted: false,
				},
				UserName: "Admin",
				Password: "hashed_password",
			},
			saveUserError:     nil,
			saveUserCallCount: 1,
			headerCallCount:   1,
			status:            http.StatusOK,
			response:          gin.H{"info": "New user successfully created"},
			secret:            "secret",
		},
		{
			name:              "Error while reading request",
			getRowDataReturn:  nil,
			getRowDataError:   errors.New("error while reading request"),
			saveUserReturn:    storage.User{},
			saveUserError:     nil,
			saveUserCallCount: 0,
			headerCallCount:   0,
			status:            http.StatusBadRequest,
			response:          gin.H{"error": "Error while reading request"},
			secret:            "secret",
		},
		{
			name:              "Error while marshalling json",
			getRowDataReturn:  []byte("WRONG JSON STRING"),
			getRowDataError:   nil,
			saveUserReturn:    storage.User{},
			saveUserError:     nil,
			saveUserCallCount: 0,
			headerCallCount:   0,
			status:            http.StatusBadRequest,
			response:          gin.H{"error": "Error while marshalling json"},
			secret:            "secret",
		},
		{
			name:              "User already exists",
			getRowDataReturn:  []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:   nil,
			saveUserReturn:    storage.User{},
			saveUserError:     errors.New("user already exists"),
			saveUserCallCount: 1,
			headerCallCount:   0,
			status:            http.StatusConflict,
			response:          gin.H{"error": "User already exists"},
		},
		{
			name:              "Internal Server Error",
			getRowDataReturn:  []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:   nil,
			saveUserReturn:    storage.User{},
			saveUserError:     errors.New("internal Server Error"),
			saveUserCallCount: 1,
			headerCallCount:   0,
			status:            http.StatusInternalServerError,
			response:          gin.H{"error": "Internal Server Error"},
		},
		{
			name:              "Create JWT Error",
			getRowDataReturn:  []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:   nil,
			saveUserReturn:    storage.User{UserName: "Admin", Password: "hashed_password"},
			saveUserError:     nil,
			saveUserCallCount: 1,
			headerCallCount:   0,
			status:            http.StatusInternalServerError,
			response:          gin.H{"error": "failed to create JWT token"},
			secret:            "",
		},
	}
	userService := mocks.NewMockUserService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		requestContext.EXPECT().GetRawData().Return(tt.getRowDataReturn, tt.getRowDataError)
		requestContext.EXPECT().Header(gomock.Any(), gomock.Any()).Times(tt.headerCallCount)
		requestContext.EXPECT().JSON(tt.status, tt.response)

		userService.EXPECT().SaveUser(gomock.Any()).
			Return(tt.saveUserReturn, tt.saveUserError).
			Times(tt.saveUserCallCount)
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				userService: userService,
				secret:      tt.secret,
			}
			h.RegisterUser(requestContext)
		})
	}
}

func TestHandler_LoginUser(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name             string
		secret           string
		getRowDataReturn []byte
		getRowDataError  error
		getUserReturn    storage.User
		getUserError     error
		getUserCallCount int
		headerCallCount  int
		status           int
		response         gin.H
	}{
		{
			name:             "Success",
			secret:           "secret",
			getRowDataReturn: []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:  nil,
			getUserReturn: storage.User{
				Entity: storage.Entity{
					Model:     gorm.Model{ID: 1},
					IsDeleted: false,
				},
				UserName: "Admin",
				Password: "hashed_password",
			},
			getUserError:     nil,
			getUserCallCount: 1,
			headerCallCount:  1,
			status:           http.StatusOK,
			response:         gin.H{"info": "login successful"},
		},
		{
			name:             "Error while reading request",
			getRowDataReturn: nil,
			getRowDataError:  errors.New("error while reading request"),
			getUserReturn:    storage.User{},
			getUserError:     nil,
			getUserCallCount: 0,
			headerCallCount:  0,
			status:           http.StatusBadRequest,
			response:         gin.H{"error": "Error while reading request"},
			secret:           "secret",
		},
		{
			name:             "Error while marshalling json",
			getRowDataReturn: []byte("WRONG JSON STRING"),
			getRowDataError:  nil,
			getUserReturn:    storage.User{},
			getUserError:     nil,
			getUserCallCount: 0,
			headerCallCount:  0,
			status:           http.StatusBadRequest,
			response:         gin.H{"error": "Error while marshalling json"},
			secret:           "secret",
		},
		{
			name:             "Wrong password",
			getRowDataReturn: []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:  nil,
			getUserReturn:    storage.User{},
			getUserError:     errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
			getUserCallCount: 1,
			headerCallCount:  0,
			status:           http.StatusUnauthorized,
			response:         gin.H{"error": "wrong password"},
		},
		{
			name:             "Internal Server Error",
			getRowDataReturn: []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:  nil,
			getUserReturn:    storage.User{},
			getUserError:     errors.New("internal Server Error"),
			getUserCallCount: 1,
			headerCallCount:  0,
			status:           http.StatusInternalServerError,
			response:         gin.H{"error": "failed to get user by username"},
		},
		{
			name:             "Create JWT Error",
			getRowDataReturn: []byte("{\n\t\"login\": \"Admin\",\n\t\"password\": \"<password>\"\n}"),
			getRowDataError:  nil,
			getUserReturn:    storage.User{UserName: "Admin", Password: "hashed_password"},
			getUserError:     nil,
			getUserCallCount: 1,
			headerCallCount:  0,
			status:           http.StatusInternalServerError,
			response:         gin.H{"error": "failed to create JWT token"},
			secret:           "",
		},
	}
	userService := mocks.NewMockUserService(ctrl)
	requestContext := mocks.NewMockRequestContext(ctrl)
	for _, tt := range tests {
		requestContext.EXPECT().GetRawData().Return(tt.getRowDataReturn, tt.getRowDataError)
		requestContext.EXPECT().Header(gomock.Any(), gomock.Any()).Times(tt.headerCallCount)
		requestContext.EXPECT().JSON(tt.status, tt.response)

		userService.EXPECT().GetUserByUserName(gomock.Any()).
			Return(tt.getUserReturn, tt.getUserError).
			Times(tt.getUserCallCount)

		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				userService: userService,
				secret:      tt.secret,
			}
			h.LoginUser(requestContext)
		})
	}
}
