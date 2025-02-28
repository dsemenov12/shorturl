package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString(t *testing.T) {
	// Тестирование успешного создания JWT токена
	t.Run("Success", func(t *testing.T) {
		userID := "user123"

		token, err := BuildJWTString(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Пытаемся разобрать токен
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		// Проверяем claims
		claims, ok := parsedToken.Claims.(*Claims)
		assert.True(t, ok)
		assert.Equal(t, userID, claims.UserID)
	})

	// Тестирование ошибки при создании токена
	t.Run("Error", func(t *testing.T) {
		// Тут можно проверить различные случаи, например, когда неправильный секретный ключ
		// В данном случае ошибки при создании токена быть не должно, так как SecretKey не меняется
		// Тест только на успешный кейс
	})
}

func TestGetUserID(t *testing.T) {
	// Тестирование извлечения UserID из валидного JWT токена
	t.Run("Success", func(t *testing.T) {
		userID := "user123"
		token, err := BuildJWTString(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Извлекаем UserID из токена
		extractedUserID, err := GetUserID(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})

	// Тестирование ошибки при неверном токене
	t.Run("InvalidToken", func(t *testing.T) {
		invalidToken := "invalid_token"

		_, err := GetUserID(invalidToken)
		assert.Error(t, err)
	})

	// Тестирование ошибки при истекшем токене
	t.Run("ExpiredToken", func(t *testing.T) {
		// Создаем токен с истекшим сроком действия
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Токен просрочен
			},
			UserID: "user123",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(SecretKey))
		assert.NoError(t, err)

		// Пытаемся извлечь UserID из просроченного токена
		_, err = GetUserID(tokenString)
		assert.Error(t, err)
	})
}

func TestTokenExpiry(t *testing.T) {
	// Тестирование срока действия токена
	t.Run("TokenExpiry", func(t *testing.T) {
		userID := "user123"
		token, err := BuildJWTString(userID)
		assert.NoError(t, err)

		// Пытаемся разобрать токен и проверяем срок действия
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(*Claims)
		assert.True(t, ok)
		// Проверка на то, что срок действия токена правильный (токен должен быть валиден 24 часа)
		assert.WithinDuration(t, claims.ExpiresAt.Time, time.Now().Add(TokenExp), time.Minute)
	})
}
