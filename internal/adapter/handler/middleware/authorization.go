package middleware

import (
	"context"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/utils"
	"github.com/bubaew95/yandex-diploma/pkg/token"
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

	return "", apperrors.TokenNotFoundErr
}

func AuthMiddleware(cfg *conf.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := getToken(r)
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, response.Response{
					Status:  "failed",
					Message: "Unauthorized",
				})
				return
			}

			jwtToken := token.NewJwtToken(cfg.SecretKey)
			user, err := jwtToken.EncodeToken(tokenStr)
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, response.Response{
					Status:  "failed",
					Message: "Unauthorized",
				})
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			request := r.WithContext(ctx)

			next.ServeHTTP(w, request)
		})
	}
}
