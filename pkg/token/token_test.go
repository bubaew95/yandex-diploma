package token

import (
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJwtToken_GenerateToken(t *testing.T) {
	t.Parallel()

	user := userentity.User{
		ID:    1,
		Login: "test",
	}

	jwtToken := NewJwtToken("123test")
	token, err := jwtToken.GenerateToken(user)
	require.NoError(t, err)

	require.NotEmpty(t, token)
}

func TestJwtToken_ParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMxMjA5NjMsIklEIjoxLCJMb2dpbiI6InRlc3QifQ.ERt9RdRr_Kx7x8B6NI-tvj9u9kO_lBNPHEU18f2qtKo"

	jwtToken := NewJwtToken("123test")
	user, err := jwtToken.EncodeToken(token)
	require.NoError(t, err)

	userEntity := userentity.User{
		ID:    1,
		Login: "test",
	}

	assert.Equal(t, userEntity, user)
}
