package authhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/dsemenov12/shorturl/internal/auth"
	"github.com/google/uuid"
)

// AuthHandle является middleware-функцией для обработки авторизации пользователей.
// Функция проверяет наличие JWT-токена в cookie запроса, если токен отсутствует, генерирует новый токен и устанавливает его в cookie.
// Если токен присутствует, извлекает идентификатор пользователя из токена и добавляет его в контекст запроса.
//
// handlerFunc: Функция, которая будет вызвана после успешной авторизации пользователя.
//
// Возвращаемое значение: возвращает обработчик HTTP-запроса (http.HandlerFunc), который выполняет проверку авторизации
// и добавляет идентификатор пользователя в контекст запроса перед вызовом основной логики.
func AuthHandle(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		jwtToken, err := r.Cookie("JWT")
		if err != nil {
			id := uuid.New()
			userID = id.String()

			tokenString, err := auth.BuildJWTString(userID)
			if err != nil {
				return
			}
			cookie := &http.Cookie{
				Name:    "JWT",
				Value:   tokenString,
				Expires: time.Now().Add(24 * time.Hour),
				Path:    "/",
			}

			http.SetCookie(w, cookie)
		} else {
			userID, err = auth.GetUserID(jwtToken.Value)
			if err != nil || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), auth.UserIDKey, userID))

		handlerFunc(w, r)
	})
}
