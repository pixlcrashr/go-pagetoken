package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand/v2"
)

func randKey(size int) ([]byte, error) {
	key := make([]byte, size)
	r := rand.ChaCha8{}
	_, err := r.Read(key)
	return key, err
}

func Rand16ByteKey() ([]byte, error) {
	return randKey(16)
}

func Rand24ByteKey() ([]byte, error) {
	return randKey(24)
}

func Rand32ByteKey() ([]byte, error) {
	return randKey(32)
}

type Encrypter interface {
	Encrypt(d []byte) (string, error)
}

type Decrypter interface {
	Decrypt(token string) ([]byte, error)
}

type Crypter interface {
	Encrypter
	Decrypter
}

type AEADEncryptor struct {
	aead cipher.AEAD
}

func NewAEADEncryptor(key []byte) (*AEADEncryptor, error) {
	// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &AEADEncryptor{aead: aead}, nil
}

func (e *AEADEncryptor) Encrypt(d []byte) (string, error) {
	nonce := make([]byte, e.aead.NonceSize())
	r := rand.ChaCha8{}
	if _, err := r.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and append nonce to ciphertext
	// Layout: nonce || ciphertext || tag
	ciphertext := e.aead.Seal(nonce, nonce, d, nil)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *AEADEncryptor) Decrypt(token string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	nonceSize := e.aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract nonce and decrypt
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := e.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
