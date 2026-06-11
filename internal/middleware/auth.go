package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const AccountIDKey contextKey = "account_id"

// DATA STRUCTURE CONCEPT:
// Context acts like a key-value map attached to request lifecycle.
// It allows safe metadata propagation without global state.

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["type"] != "access" {
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		accountIDRaw, ok := claims["account_id"]
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		accountIDStr, ok := accountIDRaw.(string)
		if !ok {
			http.Error(w, "Invalid account ID format", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AccountIDKey, accountIDStr)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
}
