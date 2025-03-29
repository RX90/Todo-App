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
	"github.com/stretchr/testify/assert"
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
			inputBody: `{"username":"Username-123", "password":"Password_123"}`,
			inputUser: todo.User{
				Username: "Username-123",
				Password: "Password_123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid input",
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
			name:                 "Empty structure",
			inputBody:            `{}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't validate input: username is empty"}`,
		},
		{
			name:                 "Empty required field",
			inputBody:            `{"username":"", "password":"hello world"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't validate input: username is empty"}`,
		},
		{
			name:      "Password is short",
			inputBody: `{"username":"Michael", "password":"7symbol"}`,
			inputUser: todo.User{
				Username: "Michael",
				Password: "7symbol",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't validate input: password is less than 8 characters"}`,
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

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
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
			inputBody: `{"username":"Existing_User", "password":"Correct-Password1"}`,
			inputUser: todo.User{
				Username: "Existing_User",
				Password: "Correct-Password1",
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
			inputBody: `{"username":"Non-existing-User", "password":"Incorrect-Password1"}`,
			inputUser: todo.User{
				Username: "Non-existing-User",
				Password: "Incorrect-Password1",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().GetUserId(user).Return("", errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"err":"can't get user id: user not found"}`,
		},
		{
			name:      "Create access token error",
			inputBody: `{"username":"Existing-User", "password":"Correct-Password1"}`,
			inputUser: todo.User{
				Username: "Existing-User",
				Password: "Correct-Password1",
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
			inputBody: `{"username":"Existing-User", "password":"Correct-Password1"}`,
			inputUser: todo.User{
				Username: "Existing-User",
				Password: "Correct-Password1",
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

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			if testCase.expectedRefreshToken != "" {
				assert.Equal(t, testCase.expectedRefreshToken, w.Result().Cookies()[0].Value)
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
			name:         "Valid input 1",
			accessToken:  "expired-valid-access-token",
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
			name:         "Valid input 2",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			mockBehavior: func(s *mock_service.MockAuthorization, accessToken, refreshToken string) {
				s.EXPECT().ParseAccessToken(accessToken).Return("1", nil)
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

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			if testCase.expectedRefreshToken != "" {
				assert.Equal(t, testCase.expectedRefreshToken, w.Result().Cookies()[0].Value)
			}
		})
	}
}

func TestAuth_logout(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid input",
			mockBehavior: func(s *mock_service.MockAuthorization) {
				s.EXPECT().DeleteRefreshToken("1").Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name: "Service error",
			mockBehavior: func(s *mock_service.MockAuthorization) {
				s.EXPECT().DeleteRefreshToken("1").Return(errors.New("db error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"error occured on deleting refresh token: db error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/logout", func(ctx *gin.Context) {
				ctx.Set(userCtx, "1")
				ctx.Next()
			}, handler.logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/logout", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
