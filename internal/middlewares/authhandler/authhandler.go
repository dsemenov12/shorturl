package authhandler

import (
	"net/http"
	"time"
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TokenExp = time.Hour * 24
const SecretKey = "supersecretkey"

type Claims struct {
    jwt.RegisteredClaims
    UserID string
}

func AuthHandle(handlerFunc http.HandlerFunc) http.HandlerFunc {
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

		handlerFunc(w, r)
	})
}

func buildJWTString(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims {
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

func getUserID(tokenString string) (string, error) {
    claims := &Claims{}

    _, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
        return []byte(SecretKey), nil
    })
	if err != nil {
		return "", err
	}

    return claims.UserID, nil
}