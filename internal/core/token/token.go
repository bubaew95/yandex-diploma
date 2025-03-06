package token

import (
	"github.com/bubaew95/yandex-diploma/internal/core/entity"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TOKEN_EXP = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	entity.User
}

type JwtToken struct {
	secretKey string
}

func NewJwtToken(secretKey string) *JwtToken {
	return &JwtToken{
		secretKey: secretKey,
	}
}

func (j JwtToken) GenerateToken(user entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		User: user,
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j JwtToken) EncodeToken(tokenString string) (user entity.User, err error) {
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return entity.User{}, err
	}

	return claims.User, nil
}
