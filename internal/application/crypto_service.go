package application

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"golang.org/x/crypto/bcrypt"
)

type CryptoService interface {
	GenerateHash(text string) (string, error)
	CompareHash(hash, text string) error
	Encrypt(text string, key enum.CryptoKey) (string, error)
	Decrypt(text string, key enum.CryptoKey) (string, error)
}

type cryptoService struct {
	c    *config.Config
	keys map[enum.CryptoKey]string
}

func NewCryptoService(c *config.Config) CryptoService {
	keys := make(map[enum.CryptoKey]string)

	keys[enum.CryptoBrokerKey] = c.Crytpo.BrokersAESKey

	return &cryptoService{c, keys}
}

func (s *cryptoService) GenerateHash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), int(s.c.Crytpo.HashCost))
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *cryptoService) CompareHash(hash, text string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
}

func (s *cryptoService) Encrypt(text string, key enum.CryptoKey) (string, error) {
	keyBytes, err := hex.DecodeString(s.keys[key])
	if err != nil {
		return "", err
	}

	textBytes := []byte(text)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipText := gcm.Seal(nonce, nonce, textBytes, nil)

	return hex.EncodeToString(cipText), nil
}

func (s *cryptoService) Decrypt(text string, key enum.CryptoKey) (string, error) {
	keyBytes, err := hex.DecodeString(s.keys[key])
	if err != nil {
		return "", err
	}

	textBytes, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	nonce, cipText := textBytes[:nonceSize], textBytes[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
