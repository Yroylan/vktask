package middleware

import (
	"context"
	"net/http"
	"strings"
	"vtask/handlers"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return handlers.JwtKey, nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


