# cmd/gophermart

В данной директории будет содержаться код накопительной системы лояльности, который скомпилируется в бинарное
приложение.

package middleware

import (
"errors"
"net/http"
"strings"

    "github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecretkey"

// Парсим токен из Cookie или Headers
func getToken(r *http.Request) (string, error) {
// 1. Проверяем заголовок Authorization: Bearer <token>
authHeader := r.Header.Get("Authorization")
if authHeader != "" {
parts := strings.Split(authHeader, " ")
if len(parts) == 2 && parts[0] == "Bearer" {
return parts[1], nil
}
}

    // 2. Проверяем Cookie
    cookie, err := r.Cookie("auth_token")
    if err == nil {
        return cookie.Value, nil
    }

    return "", errors.New("token not found")
}

// Middleware для проверки токена
func AuthMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
tokenStr, err := getToken(r)
if err != nil {
http.Error(w, "Unauthorized", http.StatusUnauthorized)
return
}

        // Разбираем токен
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}