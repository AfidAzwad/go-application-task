package middleware

import (
	"context"
	"fmt"
	"go-application-task/pkg/utils"
	"net/http"
	"os"
	"strings"
)

// JWTMiddleware is the middleware that validates the JWT token
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No change after trim
			http.Error(w, "Missing Bearer Token", http.StatusUnauthorized)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			fmt.Println("Error validating token:", err) // Add more logs for debugging
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Valid token, claims:", claims) // Ensure this log is printed

		ctx := r.Context()
		ctx = context.WithValue(ctx, "email", claims.Email)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
