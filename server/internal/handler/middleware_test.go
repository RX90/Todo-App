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
	"github.com/magiconair/properties/assert"
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

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestMiddlware_inputValidate(t *testing.T) {
	testTable := []struct {
		name          string
		input         []string
		expectedError error
	}{
		{
			name:          "Valid input",
			input:         []string{"Test Username ☻", "Test Password"},
			expectedError: nil,
		},
		{
			name:          "Username length exceeds",
			input:         []string{"Test username for testing the inputValidate function", "Test Password"},
			expectedError: errors.New("input exceeds 32 characters"),
		},
		{
			name:          "No Password",
			input:         []string{""},
			expectedError: errors.New("empty input"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, inputValidate(testCase.input...), testCase.expectedError)
		})
	}
}
