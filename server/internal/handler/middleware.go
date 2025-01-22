package handler

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
	userCtx    = "userId"
)

var ok = map[rune]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
	' ': true,

	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true,
	'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true,
	'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true,
	'v': true, 'w': true, 'x': true, 'y': true, 'z': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true,
	'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true,
	'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true,
	'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,

	'а': true, 'б': true, 'в': true, 'г': true, 'д': true, 'е': true, 'ё': true,
	'ж': true, 'з': true, 'и': true, 'й': true, 'к': true, 'л': true, 'м': true,
	'н': true, 'о': true, 'п': true, 'р': true, 'с': true, 'т': true, 'у': true,
	'ф': true, 'х': true, 'ц': true, 'ч': true, 'ш': true, 'щ': true, 'ъ': true,
	'ы': true, 'ь': true, 'э': true, 'ю': true, 'я': true,
	'А': true, 'Б': true, 'В': true, 'Г': true, 'Д': true, 'Е': true, 'Ё': true,
	'Ж': true, 'З': true, 'И': true, 'Й': true, 'К': true, 'Л': true, 'М': true,
	'Н': true, 'О': true, 'П': true, 'Р': true, 'С': true, 'Т': true, 'У': true,
	'Ф': true, 'Х': true, 'Ц': true, 'Ч': true, 'Ш': true, 'Щ': true, 'Ъ': true,
	'Ы': true, 'Ь': true, 'Э': true, 'Ю': true, 'Я': true,
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "access token is empty"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "access token is invalid"})
		return
	}

	userId, err := h.services.Authorization.ParseAccessToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	c.Set(userCtx, userId)
	c.Next()
}

func getUserCtx(c *gin.Context) (string, error) {
	idAny, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	id, ok := idAny.(string)
	if !ok {
		return "", errors.New("user id not found")
	}

	return id, nil
}

func inputValidate(input ...string) error {
	for _, data := range input {
		if utf8.RuneCountInString(data) > 32 {
			return errors.New("input exceeds 32 characters")
		}
		for _, char := range data {
			if !ok[char] {
				return errors.New("input has invalid character(s)")
			}
		}
	}

	return nil
}
