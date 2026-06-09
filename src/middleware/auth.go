package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/kaanchinar/todo-app/config"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid subject claim")
				return
			}

			userID, err := strconv.Atoi(sub)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid subject claim")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserIDKey).(int)
	return id, ok
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + msg + `"}`))
}
