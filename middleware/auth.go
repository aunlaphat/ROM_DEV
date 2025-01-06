package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type contextKey string

const (
	ContextUserRole contextKey = "userRole"
)

func AuthMiddleware(logger *zap.Logger, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Return the secret key
				return []byte("your-secret-key"), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userRoles, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "User roles not found in token", http.StatusUnauthorized)
				return
			}

			// Split the roles into a slice
			roles := strings.Split(userRoles, ",")

			// Check if the user has any of the required roles
			for _, requiredRole := range requiredRoles {
				for _, role := range roles {
					if role == requiredRole {
						ctx := context.WithValue(r.Context(), ContextUserRole, userRoles)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
