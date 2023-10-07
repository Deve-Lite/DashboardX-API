package application

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type CryptoService interface {
	Encrypt(text string, key string) (string, error)
	Decrypt(text string, key string) (string, error)
}

type cryptoService struct{}

func NewCryptoService() CryptoService {
	return &cryptoService{}
}

func (*cryptoService) Encrypt(text string, key string) (string, error) {
	keyBytes, err := hex.DecodeString(key)
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

func (*cryptoService) Decrypt(text string, key string) (string, error) {
	keyBytes, err := hex.DecodeString(key)
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
