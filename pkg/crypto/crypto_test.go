package crypto

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrypto_Encode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		correct bool
	}{
		{
			name:    "Correct encode",
			data:    "28bc66e534463037a2b39be40879cc783076219b",
			correct: true,
		},
		{
			name:    "Incorrect encode",
			data:    "28bc66e53446",
			correct: false,
		},
	}

	crypto := NewCrypto("123test")
	cryptoEncode, err := crypto.Encode("test")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.correct {
				assert.Equal(t, tt.data, cryptoEncode)
			} else {
				assert.NotEqual(t, tt.data, cryptoEncode)
			}
		})
	}

}

func TestCrypto_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		correct bool
	}{
		{
			name:    "Correct decode",
			data:    "test",
			correct: true,
		},
		{
			name:    "Incorrect decode",
			data:    "test_1",
			correct: false,
		},
	}

	crypto := NewCrypto("123test")
	cryptoDecode, err := crypto.Decode("28bc66e534463037a2b39be40879cc783076219b")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.correct {
				assert.Equal(t, tt.data, cryptoDecode)
			} else {
				assert.NotEqual(t, tt.data, cryptoDecode)
			}
		})
	}
}
