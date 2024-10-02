package service_test

import (
	"bookstore/internal/model"
	"bookstore/internal/service"
	"bookstore/pkg/utils"
	"bookstore/test/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCustomerService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerRepository(ctrl)
	service := service.NewCustomerService(mockRepo)

	customer := &model.Customer{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mockRepo.EXPECT().Register(customer).Return(nil).Times(1)

	err := service.Register(customer)
	assert.NoError(t, err)
}

func TestCustomerService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerRepository(ctrl)
	service := service.NewCustomerService(mockRepo)

	email := "test@example.com"
	password := "password"

	hashedPassword, _ := utils.HashPassword(password)
	customer := &model.Customer{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}

	mockRepo.EXPECT().Login(email, password).Return(customer, nil).Times(1)

	token, err := service.Login(email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestCustomerService_Login_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerRepository(ctrl)
	service := service.NewCustomerService(mockRepo)

	email := "test@example.com"
	password := "wrongpassword"

	mockRepo.EXPECT().Login(email, password).Return(nil, errors.New("invalid credentials")).Times(1)

	token, err := service.Login(email, password)
	assert.Error(t, err)
	assert.Empty(t, token)
}
