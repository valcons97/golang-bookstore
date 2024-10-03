package handler_test

import (
	"bookstore/internal/handler"
	"bookstore/internal/handler/request"
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"bookstore/test/mocks"
	"bytes"
	"encoding/json"
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
		token := "testToken"

		mockCustomerService.EXPECT().
			Login("test@example.com", "password").
			Return(token, nil)

		loginRequest := request.LoginRequest{Email: "test@example.com", Password: "password"}
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

	})

	t.Run("invalid email or password", func(t *testing.T) {
		mockCustomerService.EXPECT().
			Login("wrong@example.com", "wrongpassword").
			Return("", utils.ErrWrongPassword)

		// Prepare request with incorrect credentials
		loginRequest := request.LoginRequest{Email: "wrong@example.com", Password: "wrongpassword"}
		jsonReq, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

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
