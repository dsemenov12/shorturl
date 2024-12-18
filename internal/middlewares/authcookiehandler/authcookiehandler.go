package authcookiehandler

import (
	"context"
	"net/http"

	"github.com/dsemenov12/shorturl/internal/auth"
)

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