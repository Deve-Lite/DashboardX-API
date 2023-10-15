package application_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/go-playground/assert"
)

func generateAESKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	return hex.EncodeToString(bytes)
}

func TestCryptoService(t *testing.T) {
	cs := application.NewCryptoService()
	key := generateAESKey()

	t.Run("should encrypt provided value", func(t *testing.T) {
		s := "text to encrypt"
		r, err := cs.Encrypt(s, key)

		assert.NotEqual(t, s, r)
		assert.Equal(t, err, nil)
	})

	t.Run("should decrypt previously encrypted value ", func(t *testing.T) {
		s := "text to encrypt"
		r, err := cs.Encrypt(s, key)
		assert.NotEqual(t, s, r)
		assert.Equal(t, err, nil)

		r2, err2 := cs.Decrypt(r, key)
		assert.Equal(t, s, r2)
		assert.Equal(t, err2, nil)
	})

	t.Run("should raise error when encryption key is invalid", func(t *testing.T) {
		r, err := cs.Encrypt("text", "invalid-aes-key")
		assert.Equal(t, r, "")
		assert.NotEqual(t, err, nil)
	})

	t.Run("should raise error when decryption key is invalid", func(t *testing.T) {
		s, _ := cs.Encrypt("text", key)

		r, err := cs.Decrypt(s, "invalid-aes-key")
		assert.Equal(t, r, "")
		assert.NotEqual(t, err, nil)
	})

	t.Run("should raise error encrypted value is invalid", func(t *testing.T) {
		r, err := cs.Decrypt("not-really-encrypted", key)
		assert.Equal(t, r, "")
		assert.NotEqual(t, err, nil)
	})
}
