// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/keyjin88/go-loyalty-system/internal/app/handlers (interfaces: WithdrawService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/keyjin88/go-loyalty-system/internal/app/storage"
)

// MockWithdrawService is a mock of WithdrawService interface.
type MockWithdrawService struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawServiceMockRecorder
}

// MockWithdrawServiceMockRecorder is the mock recorder for MockWithdrawService.
type MockWithdrawServiceMockRecorder struct {
	mock *MockWithdrawService
}

// NewMockWithdrawService creates a new mock instance.
func NewMockWithdrawService(ctrl *gomock.Controller) *MockWithdrawService {
	mock := &MockWithdrawService{ctrl: ctrl}
	mock.recorder = &MockWithdrawServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawService) EXPECT() *MockWithdrawServiceMockRecorder {
	return m.recorder
}

// GetAllWithdrawals mocks base method.
func (m *MockWithdrawService) GetAllWithdrawals(arg0 uint) ([]storage.WithdrawResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWithdrawals", arg0)
	ret0, _ := ret[0].([]storage.WithdrawResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWithdrawals indicates an expected call of GetAllWithdrawals.
func (mr *MockWithdrawServiceMockRecorder) GetAllWithdrawals(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWithdrawals", reflect.TypeOf((*MockWithdrawService)(nil).GetAllWithdrawals), arg0)
}

// SaveWithdraw mocks base method.
func (m *MockWithdrawService) SaveWithdraw(arg0 storage.WithdrawRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWithdraw", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveWithdraw indicates an expected call of SaveWithdraw.
func (mr *MockWithdrawServiceMockRecorder) SaveWithdraw(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWithdraw", reflect.TypeOf((*MockWithdrawService)(nil).SaveWithdraw), arg0)
}
