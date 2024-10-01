package handler_test

import (
	"bookstore/internal/handler"
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"bookstore/test/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCustomerHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCustomerService := mocks.NewMockCustomerService(ctrl)
	router := gin.Default()

	customerHandler := handler.NewCustomerHandler(mockCustomerService)
	router.POST("/login", customerHandler.Login)

	t.Run("success", func(t *testing.T) {
		mockCustomer := model.Customer{
			ID:      int64(1),
			Email:   "test@example.com",
			Name:    "Test User",
			Address: "123 Test Street",
		}
		mockCustomerService.EXPECT().
			Login("test@example.com", "password").
			Return(&mockCustomer, nil)

		token, _ := utils.GenerateToken(mockCustomer.ID, mockCustomer.Email)

		loginRequest := handler.LoginRequest{Email: "test@example.com", Password: "password"}
		jsonReq, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, "Login successful", actualResponse["message"])
		assert.Equal(t, token, actualResponse["token"])

		customer := actualResponse["customer"].(map[string]interface{})
		assert.Equal(t, mockCustomer.ID, int64(customer["id"].(float64)))
		assert.Equal(t, mockCustomer.Email, customer["email"])
		assert.Equal(t, mockCustomer.Name, customer["name"])
		assert.Equal(t, mockCustomer.Address, customer["address"])
	})

	t.Run("invalid email or password", func(t *testing.T) {
		mockCustomerService.EXPECT().
			Login("wrong@example.com", "wrongpassword").
			Return(&model.Customer{}, errors.New("invalid credentials"))

		// Prepare request with incorrect credentials
		loginRequest := handler.LoginRequest{Email: "wrong@example.com", Password: "wrongpassword"}
		jsonReq, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCustomerHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCustomerService := mocks.NewMockCustomerService(ctrl)
	router := gin.Default()

	customerHandler := handler.NewCustomerHandler(mockCustomerService)
	router.POST("/register", customerHandler.Register)

	t.Run("success", func(t *testing.T) {
		newCustomer := model.Customer{
			Email:    "newuser@example.com",
			Name:     "New User",
			Address:  "456 New Street",
			Password: "hashedpassword",
		}

		mockCustomerService.EXPECT().Register(gomock.Any()).Return(nil)

		jsonReq, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		expectedResponse := gin.H{"message": "Customer registered successfully"}
		var actualResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBuffer([]byte("{bad json}")),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
