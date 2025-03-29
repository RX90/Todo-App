package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RX90/Todo-App/server/internal/service"
	mock_service "github.com/RX90/Todo-App/server/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Valid header",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return("1", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "1",
		},
		{
			name:                 "No header",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"auth header is empty"}`,
		},
		{
			name:                 "Invalid header 1",
			headerName:           "Authorization",
			headerValue:          "Bearertoken",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"auth header is invalid"}`,
		},
		{
			name:                 "Invalid header 2",
			headerName:           "Authorization",
			headerValue:          "Berer token",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"auth header is invalid"}`,
		},
		{
			name:                 "Invalid header 3",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"auth header is invalid"}`,
		},
		{
			name:        "Service error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseAccessToken(token).Return("", errors.New("token has expired"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"token has expired"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get(userCtx)
				ctx.String(http.StatusOK, id.(string))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestMiddleware_authInputValidation(t *testing.T) {
	testTable := []struct {
		name          string
		username      string
		password      string
		expectedError error
	}{
		{
			name:          "Valid input",
			username:      "New-Username",
			password:      "New_Password123",
			expectedError: nil,
		},
		{
			name:          "Empty username",
			username:      "",
			password:      "New_Password123",
			expectedError: errors.New("username is empty"),
		},
		{
			name:          "Empty password",
			username:      "New-Username",
			password:      "",
			expectedError: errors.New("password is empty"),
		},
		{
			name:          "Long username",
			username:      "Картофелекопатель",
			password:      "New_Password123",
			expectedError: errors.New("username exceeds 32 bytes"),
		},
		{
			name:          "Long Password",
			username:      "New-Username",
			password:      "New-Very-Very-Very-Very-Big-Password",
			expectedError: errors.New("password exceeds 32 bytes"),
		},
		{
			name:          "Short username",
			username:      "1",
			password:      "New_Password123",
			expectedError: errors.New("username is less than 3 characters"),
		},
		{
			name:          "Short password",
			username:      "New-Username",
			password:      "qwerty",
			expectedError: errors.New("password is less than 8 characters"),
		},
		{
			name:          "Invalid characters in username",
			username:      "Валерий",
			password:      "New_Password123",
			expectedError: errors.New("username has invalid character"),
		},
		{
			name:          "Invalid characters in password",
			username:      "New-Username",
			password:      "qwerty 123",
			expectedError: errors.New("password has invalid character"),
		},
		{
			name:          "No letter(s) in password",
			username:      "New-Username",
			password:      "1234567890",
			expectedError: errors.New("password must contain at least one english letter and one digit"),
		},
		{
			name:          "No digit(s) in password",
			username:      "New-Username",
			password:      "New_Password",
			expectedError: errors.New("password must contain at least one english letter and one digit"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedError, authInputValidation(testCase.username, testCase.password))
		})
	}
}
