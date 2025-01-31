package service

import (
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func newAccessToken(userId string) string {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			Subject:   userId,
		},
		userId,
	})

	tokenString, err := accessToken.SignedString([]byte(signingKey))
	if err != nil {
		return ""
	}

	return tokenString
}

func TestAuth_ParseAccessToken(t *testing.T) {

	testTable := []struct {
		name           string
		token          string
		expectedUserId string
		expectedError  error
	}{
		{
			name:           "Valid access token",
			token:          newAccessToken("1"),
			expectedUserId: "1",
			expectedError:  nil,
		},
		{
			name:           "Expired token",
			token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY2ODU1MTgsInN1YiI6IjEiLCJ1c2VyX2lkIjoiMSJ9.33L0FPQyoCDBK28NoLIK-X0_xrpxvwJfIKkN2dpgayg",
			expectedUserId: "1",
			expectedError:  errors.New("token has expired"),
		},
		{
			name:           "Invalid token",
			token:          "invalid.token.string",
			expectedUserId: "",
			expectedError:  errors.New(`invalid character '\u008a' looking for beginning of value`),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			authService := AuthService{}

			userId, err := authService.ParseAccessToken(testCase.token)

			assert.Equal(t, testCase.expectedUserId, userId)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
