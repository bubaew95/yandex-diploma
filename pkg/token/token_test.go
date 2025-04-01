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
	userEntity := userentity.User{
		ID:    1,
		Login: "test",
	}

	jwtToken := NewJwtToken("123test")
	token, err := jwtToken.GenerateToken(userEntity)
	require.NoError(t, err)

	user, err := jwtToken.EncodeToken(token)
	require.NoError(t, err)

	assert.Equal(t, userEntity, user)
}
