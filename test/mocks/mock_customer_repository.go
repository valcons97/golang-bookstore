// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/customer_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	model "bookstore/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCustomerRepository is a mock of CustomerRepository interface.
type MockCustomerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCustomerRepositoryMockRecorder
}

// MockCustomerRepositoryMockRecorder is the mock recorder for MockCustomerRepository.
type MockCustomerRepositoryMockRecorder struct {
	mock *MockCustomerRepository
}

// NewMockCustomerRepository creates a new mock instance.
func NewMockCustomerRepository(ctrl *gomock.Controller) *MockCustomerRepository {
	mock := &MockCustomerRepository{ctrl: ctrl}
	mock.recorder = &MockCustomerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomerRepository) EXPECT() *MockCustomerRepositoryMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockCustomerRepository) Login(email, password string) (*model.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", email, password)
	ret0, _ := ret[0].(*model.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockCustomerRepositoryMockRecorder) Login(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockCustomerRepository)(nil).Login), email, password)
}

// Register mocks base method.
func (m *MockCustomerRepository) Register(customer *model.Customer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", customer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockCustomerRepositoryMockRecorder) Register(customer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockCustomerRepository)(nil).Register), customer)
}
