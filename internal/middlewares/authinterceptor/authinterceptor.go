package authinterceptor

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dsemenov12/shorturl/internal/auth"
)

// AuthUnaryInterceptor является gRPC Unary Interceptor-ом для обработки авторизации пользователей по JWT-токену из cookie.
//
// Интерцептор извлекает cookie "JWT" из входящего metadata gRPC-запроса. Если токен отсутствует или недействителен,
// генерируется новый токен с уникальным идентификатором пользователя и он устанавливается в trailing metadata как Set-Cookie.
// Полученный идентификатор пользователя добавляется в контекст запроса.
//
// Возвращаемое значение: возвращает gRPC Unary interceptor, который обеспечивает авторизацию на уровне gRPC.
func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		var jwtToken string

		// Извлекаем JWT из cookie заголовков (metadata)
		if cookieHeaders := md.Get("cookie"); len(cookieHeaders) > 0 {
			for _, cookieHeader := range cookieHeaders {
				parts := strings.Split(cookieHeader, ";")
				for _, part := range parts {
					kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
					if len(kv) == 2 && kv[0] == "JWT" {
						jwtToken = kv[1]
						break
					}
				}
			}
		}

		var userID string
		var err error
		var newToken string
		needNewToken := false

		if jwtToken == "" {
			// Генерация нового токена
			id := uuid.New()
			userID = id.String()
			newToken, err = auth.BuildJWTString(userID)
			if err != nil {
				return nil, err
			}
			needNewToken = true
		} else {
			// Проверка JWT
			userID, err = auth.GetUserID(jwtToken)
			if err != nil || userID == "" {
				// Некорректный токен — как в HTTP middleware, авторизуем заново
				id := uuid.New()
				userID = id.String()
				newToken, err = auth.BuildJWTString(userID)
				if err != nil {
					return nil, err
				}
				needNewToken = true
			}
		}

		// Добавляем userID в context
		ctx = context.WithValue(ctx, auth.UserIDKey, userID)

		// Добавляем Set-Cookie в trailing metadata (может быть перехвачено grpc-gateway)
		if needNewToken {
			setCookie := &http.Cookie{
				Name:     "JWT",
				Value:    newToken,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			trailer := metadata.Pairs("Set-Cookie", setCookie.String())
			grpc.SetTrailer(ctx, trailer)
		}

		return handler(ctx, req)
	}
}
