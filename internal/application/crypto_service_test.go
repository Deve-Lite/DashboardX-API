package application_test

import (
	"testing"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/go-playground/assert"
)

func TestCryptoService(t *testing.T) {
	c := config.NewConfig(config.GetDefaultPath("test.env"))
	cs := application.NewCryptoService(c)

	t.Run("should encrypt provided value", func(t *testing.T) {
		s := "text to encrypt"
		r, err := cs.Encrypt(s, enum.CryptoBrokerKey)

		assert.NotEqual(t, s, r)
		assert.Equal(t, err, nil)
	})

	t.Run("should decrypt previously encrypted value ", func(t *testing.T) {
		s := "text to encrypt"
		r, err := cs.Encrypt(s, enum.CryptoBrokerKey)
		assert.NotEqual(t, s, r)
		assert.Equal(t, err, nil)

		r2, err2 := cs.Decrypt(r, enum.CryptoBrokerKey)
		assert.Equal(t, s, r2)
		assert.Equal(t, err2, nil)
	})

	t.Run("should raise error when encryption key is invalid", func(t *testing.T) {
		r, err := cs.Encrypt("text", "invalid-aes-key")
		assert.Equal(t, r, "")
		assert.NotEqual(t, err, nil)
	})

	t.Run("should raise error when decryption key is invalid", func(t *testing.T) {
		s, _ := cs.Encrypt("text", enum.CryptoBrokerKey)

		r, err := cs.Decrypt(s, "invalid-aes-key")
		assert.Equal(t, r, "")
		assert.NotEqual(t, err, nil)
	})

	t.Run("should raise error encrypted value is invalid", func(t *testing.T) {
		r, err := cs.Decrypt("not-really-encrypted", enum.CryptoBrokerKey)
		assert.Equal(t, r, "")
		assert.NotEqual(t, err, nil)
	})
}
