package authhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/dsemenov12/shorturl/internal/auth"
)

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
				Name: "JWT",
				Value: tokenString,
				Expires: time.Now().Add(24 * time.Hour),
				Path: "/",
			}
		
			http.SetCookie(w, cookie)
		} else {
			userID, err = auth.GetUserID(jwtToken.Value)
			if err != nil || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		
		r = r.WithContext(context.WithValue(context.Background(), auth.UserIDKey, userID))

		handlerFunc(w, r)
	})
}