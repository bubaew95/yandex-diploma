package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Crypto struct {
	SecretKey string
}

func NewCrypto(secretKey string) *Crypto {
	return &Crypto{
		SecretKey: secretKey,
	}
}

func (c Crypto) Decode(data string) (string, error) {
	aesgcm, nonce, err := aesGcm(c.SecretKey)
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (c Crypto) Encode(data string) (string, error) {
	aesgcm, nonce, err := aesGcm(c.SecretKey)
	if err != nil {
		return "", err
	}

	dst := aesgcm.Seal(nil, nonce, []byte(data), nil)
	return fmt.Sprintf("%x", dst), nil
}

func aesGcm(secretKey string) (cipher.AEAD, []byte, error) {
	key := sha256.Sum256([]byte(secretKey))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]
	return aesgcm, nonce, nil
}
