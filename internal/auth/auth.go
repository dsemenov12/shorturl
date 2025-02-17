package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthJWT описывает интерфейс для работы с JWT-токенами.
// Он предоставляет методы для создания JWT-токена и извлечения идентификатора пользователя из токена.
type AuthJWT interface {
	// BuildJWTString создает JWT-токен для указанного userID.
	// Возвращает строку с токеном или ошибку в случае неудачи.
	BuildJWTString(userID string) (string, error)

	// GetUserID извлекает идентификатор пользователя из JWT-токена.
	// Возвращает userID или ошибку в случае неудачи.
	GetUserID(tokenString string) (string, error)
}

// Claims структура представляет собой полезную нагрузку (payload) для JWT.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// TokenExp определяет время жизни JWT-токена.
const TokenExp = time.Hour * 24

// SecretKey — это секретный ключ, используемый для подписания JWT-токенов.
const SecretKey = "supersecretkey"

// userContextKey используется в контексте для хранения идентификатора пользователя.
type userContextKey string

// UserIDKey — это ключ для хранения идентификатора пользователя в контексте.
const UserIDKey userContextKey = "user_id"

// BuildJWTString создает новый JWT-токен для указанного userID.
// Этот токен подписан с использованием секретного ключа и включает в себя срок действия.
// Возвращает строку с токеном или ошибку в случае неудачи.
func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID извлекает идентификатор пользователя из JWT-токена.
// Пытается разобрать токен, используя секретный ключ и возвращает UserID, если все прошло успешно.
// Возвращает ошибку, если токен недействителен или произошла ошибка при его разборе.
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}
