package authcookiehandler

import (
	"context"
	"net/http"

	"github.com/dsemenov12/shorturl/internal/auth"
)

// AuthCookieHandle является middleware-функцией для обработки авторизации пользователей
// на основе JWT-токенов, сохраненных в cookie. Если токен действителен, извлекает идентификатор пользователя,
// добавляет его в контекст запроса и передает управление дальше в цепочку обработки.
// Если токен отсутствует или недействителен, возвращает ошибку 401 (Unauthorized).
//
// handlerFunc: Функция, которая будет вызвана после успешной авторизации пользователя.
//
// Возвращаемое значение: возвращает обработчик HTTP-запроса (http.HandlerFunc), который выполняет проверку авторизации
// перед выполнением основной логики.
func AuthCookieHandle(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		jwtToken, _ := r.Cookie("JWT")
		if jwtToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := auth.GetUserID(jwtToken.Value)
		if err != nil || userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(context.Background(), auth.UserIDKey, userID))

		handlerFunc(w, r)
	})
}
