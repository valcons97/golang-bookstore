// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/order_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	request "bookstore/internal/handler/request"
	model "bookstore/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOrderService is a mock of OrderService interface.
type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

// MockOrderServiceMockRecorder is the mock recorder for MockOrderService.
type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

// NewMockOrderService creates a new mock instance.
func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

// AddToCart mocks base method.
func (m *MockOrderService) AddToCart(customerID int, request request.AddToCartRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToCart", customerID, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToCart indicates an expected call of AddToCart.
func (mr *MockOrderServiceMockRecorder) AddToCart(customerID, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToCart", reflect.TypeOf((*MockOrderService)(nil).AddToCart), customerID, request)
}

// CreateOrderIfNotExists mocks base method.
func (m *MockOrderService) CreateOrderIfNotExists(customerID int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrderIfNotExists", customerID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrderIfNotExists indicates an expected call of CreateOrderIfNotExists.
func (mr *MockOrderServiceMockRecorder) CreateOrderIfNotExists(customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrderIfNotExists", reflect.TypeOf((*MockOrderService)(nil).CreateOrderIfNotExists), customerID)
}

// GetCart mocks base method.
func (m *MockOrderService) GetCart(customerID int) (*model.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCart", customerID)
	ret0, _ := ret[0].(*model.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCart indicates an expected call of GetCart.
func (mr *MockOrderServiceMockRecorder) GetCart(customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCart", reflect.TypeOf((*MockOrderService)(nil).GetCart), customerID)
}

// GetOrderHistory mocks base method.
func (m *MockOrderService) GetOrderHistory(customerID int, request request.HistoryRequest) ([]model.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderHistory", customerID, request)
	ret0, _ := ret[0].([]model.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderHistory indicates an expected call of GetOrderHistory.
func (mr *MockOrderServiceMockRecorder) GetOrderHistory(customerID, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderHistory", reflect.TypeOf((*MockOrderService)(nil).GetOrderHistory), customerID, request)
}

// PayOrder mocks base method.
func (m *MockOrderService) PayOrder(customerID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PayOrder", customerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PayOrder indicates an expected call of PayOrder.
func (mr *MockOrderServiceMockRecorder) PayOrder(customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PayOrder", reflect.TypeOf((*MockOrderService)(nil).PayOrder), customerID)
}

// RemoveFromCart mocks base method.
func (m *MockOrderService) RemoveFromCart(customerID, bookId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromCart", customerID, bookId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromCart indicates an expected call of RemoveFromCart.
func (mr *MockOrderServiceMockRecorder) RemoveFromCart(customerID, bookId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromCart", reflect.TypeOf((*MockOrderService)(nil).RemoveFromCart), customerID, bookId)
}
