package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RX90/Todo-App/server/internal/service"
	mock_service "github.com/RX90/Todo-App/server/internal/service/mocks"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestAuth_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid input",
			inputBody: `{"username":"Тестовый Username123", "password":"Тестовый Password123"}`,
			inputUser: todo.User{
				Username: "Тестовый Username123",
				Password: "Тестовый Password123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid input 1",
			inputBody:            `{"username":Michael, "password":true}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't bind JSON: invalid character 'M' looking for beginning of value"}`,
		},
		{
			name:                 "No required fields",
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't bind JSON: EOF"}`,
		},
		{
			name:                 "Empty required field",
			inputBody:            `{"username":"", "password":"hello world"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"empty input"}`,
		},
		{
			name:      "Invalid input 2",
			inputBody: `{"username":"Michael", "password":"7symbol"}`,
			inputUser: todo.User{
				Username: "Michael",
				Password: "7symbol",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"password is less than 8 characters"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"Username":"Test", "password":"qwerty123"}`,
			inputUser: todo.User{
				Username: "Test",
				Password: "qwerty123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create user: service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestAuth_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		expectedRefreshToken string
	}{
		{
			name:      "Valid input",
			inputBody: `{"username":"Existing User", "password":"Correct Password"}`,
			inputUser: todo.User{
				Username: "Existing User",
				Password: "Correct Password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().GetUserId(user).Return("1", nil)
				s.EXPECT().NewAccessToken("1").Return("valid-access-token", nil)
				s.EXPECT().NewRefreshToken("1").Return("valid-refresh-token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"valid-access-token"}`,
			expectedRefreshToken: "valid-refresh-token",
		},
		{
			name:      "Get user id error",
			inputBody: `{"username":"Non-existing User", "password":"Incorrect Password"}`,
			inputUser: todo.User{
				Username: "Non-existing User",
				Password: "Incorrect Password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().GetUserId(user).Return("", errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"can't get user id: user not found"}`,
		},
		{
			name:      "Create access token error",
			inputBody: `{"username":"Existing User", "password":"Correct Password"}`,
			inputUser: todo.User{
				Username: "Existing User",
				Password: "Correct Password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().GetUserId(user).Return("1", nil)
				s.EXPECT().NewAccessToken("1").Return("", errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create access token: service error"}`,
		},
		{
			name:      "Create refresh token error",
			inputBody: `{"username":"Existing User", "password":"Correct Password"}`,
			inputUser: todo.User{
				Username: "Existing User",
				Password: "Correct Password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().GetUserId(user).Return("1", nil)
				s.EXPECT().NewAccessToken("1").Return("access-token", nil)
				s.EXPECT().NewRefreshToken("1").Return("", errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create refresh token: service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
			if testCase.expectedRefreshToken != "" {
				assert.Equal(t, w.Result().Cookies()[0].Value, testCase.expectedRefreshToken)
			}
		})
	}
}

func TestAuth_refreshTokens(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, accessToken, refreshToken string)

	testTable := []struct {
		name                 string
		accessToken          string
		refreshToken         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		expectedRefreshToken string
	}{
		{
			name:         "Valid input",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {
				s.EXPECT().ParseAccessToken(accessToken).Return("1", errors.New("token has expired"))
				s.EXPECT().CheckRefreshToken("1", refreshToken).Return(nil)
				s.EXPECT().NewAccessToken("1").Return("new-valid-access-token", nil)
				s.EXPECT().NewRefreshToken("1").Return("new-valid-refresh-token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"new-valid-access-token"}`,
			expectedRefreshToken: "new-valid-refresh-token",
		},
		{
			name:                 "Header is invalid",
			accessToken:          " invalid-access-token",
			refreshToken:         "valid-refresh-token",
			mockBehavior:         func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"auth header is invalid"}`,
		},
		{
			name:                 "No access token",
			accessToken:          "",
			refreshToken:         "valid-refresh-token",
			mockBehavior:         func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"auth header is invalid"}`,
		},
		{
			name:                 "No refresh token",
			accessToken:          "valid-access-token",
			refreshToken:         "",
			mockBehavior:         func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"refresh token is missing"}`,
		},
		{
			name:         "Parse access token error",
			accessToken:  "invalid-access-token",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {
				s.EXPECT().ParseAccessToken(accessToken).Return("", errors.New("parse token error"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"access token is invalid: parse token error"}`,
		},
		{
			name:         "Check refresh token error",
			accessToken:  "valid-access-token",
			refreshToken: "invalid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {
				s.EXPECT().ParseAccessToken(accessToken).Return("1", errors.New("token has expired"))
				s.EXPECT().CheckRefreshToken("1", refreshToken).Return(errors.New("check refresh token error"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"refresh token is invalid: check refresh token error"}`,
		},
		{
			name:         "New access token error",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {
				s.EXPECT().ParseAccessToken(accessToken).Return("1", errors.New("token has expired"))
				s.EXPECT().CheckRefreshToken("1", refreshToken).Return(nil)
				s.EXPECT().NewAccessToken("1").Return("", errors.New("new access token error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create access token: new access token error"}`,
		},
		{
			name:         "New refresh token error",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {
				s.EXPECT().ParseAccessToken(accessToken).Return("1", errors.New("token has expired"))
				s.EXPECT().CheckRefreshToken("1", refreshToken).Return(nil)
				s.EXPECT().NewAccessToken("1").Return("new-valid-access-token", nil)
				s.EXPECT().NewRefreshToken("1").Return("", errors.New("new refresh token error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create refresh token: new refresh token error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.accessToken, testCase.refreshToken)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/refresh", handler.refreshTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh", nil)
			req.Header.Set(authHeader, "Bearer "+testCase.accessToken)
			req.AddCookie(&http.Cookie{Name: refresh, Value: testCase.refreshToken})

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
			if testCase.expectedRefreshToken != "" {
				assert.Equal(t, w.Result().Cookies()[0].Value, testCase.expectedRefreshToken)
			}
		})
	}
}

func TestAuth_logout(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, refreshToken string)

	testTable := []struct {
		name                 string
		refreshToken         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Valid input",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, refreshToken string) {
				s.EXPECT().DeleteRefreshToken("1", refreshToken).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "No refresh token",
			refreshToken:         "",
			mockBehavior:         func(s *mock_service.MockAuthorization, refreshToken string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"refresh token is missing"}`,
		},
		{
			name:         "Invalid refresh token",
			refreshToken: "invalid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, refreshToken string) {
				s.EXPECT().DeleteRefreshToken("1", refreshToken).Return(errors.New("delete refresh token error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't delete refresh token: delete refresh token error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.refreshToken)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/logout", func(ctx *gin.Context) {
				ctx.Set(userCtx, "1")
				ctx.Next()
			}, handler.logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/logout", nil)
			req.AddCookie(&http.Cookie{Name: refresh, Value: testCase.refreshToken})

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
