package middleware

import (
	"errors"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/token"
	"net/http"
	"strings"
)

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	cookie, err := r.Cookie("auth_token")
	if err == nil {
		return cookie.Value, nil
	}

	return "", errors.New("token not found")
}

func AuthMiddleware(cfg *conf.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := getToken(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			jwtToken := token.NewJwtToken(cfg.SecretKey)
			_, err = jwtToken.EncodeToken(tokenStr)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
