package authhandler

import (
	"net/http"
	"time"
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TOKEN_EXP = time.Hour * 24
const SECRET_KEY = "supersecretkey"

type Claims struct {
    jwt.RegisteredClaims
    UserID string
}

func AuthHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		jwtToken, err := r.Cookie("JWT")
		if err != nil {
			// установит куку
			id := uuid.New()
			userID = id.String()
			tokenString, err := buildJWTString(userID)
			if err != nil {
				return
			}
        	cookie := &http.Cookie{
				Name: "JWT",
				Value: tokenString,
				Expires: time.Now().Add(24 * time.Hour),
			}
		
			http.SetCookie(w, cookie)
		} else {
			userID, err = getUserID(jwtToken.Value)
			if err != nil || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		
		r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))

		next.ServeHTTP(w, r)
	})
}

func buildJWTString(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims {
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
        },
        UserID: userID,
    })

    tokenString, err := token.SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func getUserID(tokenString string) (string, error) {
    claims := &Claims{}

    _, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
        return []byte(SECRET_KEY), nil
    })
	if err != nil {
		return "", err
	}

    return claims.UserID, nil
}