// Package encryption provides secure encryption and decryption for page tokens.
//
// This package implements AES-GCM authenticated encryption with associated data (AEAD)
// for securing page tokens. It provides an interface-based design that allows for
// custom encryption implementations while offering a production-ready default.
//
// # Features
//
//   - AES-GCM AEAD encryption (128, 192, or 256-bit keys)
//   - ChaCha8 PRNG for nonce generation
//   - Base64 URL-safe encoding
//   - Interface-based design for extensibility
//   - Helper functions for secure key generation
//
// # Encryptor Interface
//
// The Encryptor interface allows custom encryption implementations:
//
//	type Encryptor interface {
//	    Encrypt(d []byte) (string, error)
//	    Decrypt(token string) ([]byte, error)
//	}
//
// # Example: Basic Usage
//
//	package main
//
//	import (
//	    "fmt"
//	    "log"
//
//	    "github.com/pixlcrashr/go-pagetoken/encryption"
//	)
//
//	func main() {
//	    // Generate a random 32-byte key (AES-256)
//	    key, err := encryption.Rand32ByteKey()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Create an encryptor
//	    encryptor, err := encryption.NewAEADEncryptor(key)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Encrypt some data
//	    plaintext := []byte("sensitive data")
//	    token, err := encryptor.Encrypt(plaintext)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Encrypted token: %s\n", token)
//
//	    // Decrypt the token
//	    decrypted, err := encryptor.Decrypt(token)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Decrypted data: %s\n", decrypted)
//	}
//
// # Example: Key Sizes
//
//	// AES-128 (16 bytes)
//	key128, _ := encryption.Rand16ByteKey()
//	enc128, _ := encryption.NewAEADEncryptor(key128)
//
//	// AES-192 (24 bytes)
//	key192, _ := encryption.Rand24ByteKey()
//	enc192, _ := encryption.NewAEADEncryptor(key192)
//
//	// AES-256 (32 bytes) - Recommended for production
//	key256, _ := encryption.Rand32ByteKey()
//	enc256, _ := encryption.NewAEADEncryptor(key256)
//
// # Example: Custom Encryption Implementation
//
// You can implement your own encryption strategy:
//
//	package main
//
//	import (
//	    "crypto/aes"
//	    "crypto/cipher"
//	    "encoding/base64"
//	    "errors"
//	)
//
//	type CustomEncryptor struct {
//	    block cipher.Block
//	}
//
//	func NewCustomEncryptor(key []byte) (*CustomEncryptor, error) {
//	    block, err := aes.NewCipher(key)
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &CustomEncryptor{block: block}, nil
//	}
//
//	func (e *CustomEncryptor) Encrypt(data []byte) (string, error) {
//	    // Your custom encryption logic
//	    // Must return base64-encoded string
//	    encrypted := make([]byte, len(data))
//	    // ... encryption implementation ...
//	    return base64.URLEncoding.EncodeToString(encrypted), nil
//	}
//
//	func (e *CustomEncryptor) Decrypt(token string) ([]byte, error) {
//	    ciphertext, err := base64.URLEncoding.DecodeString(token)
//	    if err != nil {
//	        return nil, err
//	    }
//	    // Your custom decryption logic
//	    decrypted := make([]byte, len(ciphertext))
//	    // ... decryption implementation ...
//	    return decrypted, nil
//	}
//
// # Example: Production Key Management
//
//	package main
//
//	import (
//	    "encoding/base64"
//	    "log"
//	    "os"
//
//	    "github.com/pixlcrashr/go-pagetoken/encryption"
//	)
//
//	func main() {
//	    // Generate a key once and store it securely
//	    key, err := encryption.Rand32ByteKey()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Store in environment (example only - use proper secrets management)
//	    keyEncoded := base64.StdEncoding.EncodeToString(key)
//	    log.Printf("Store this key securely: %s", keyEncoded)
//
//	    // Load from environment in production
//	    keyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("PAGE_TOKEN_KEY"))
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    encryptor, err := encryption.NewAEADEncryptor(keyBytes)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Use the encryptor
//	    _ = encryptor
//	}
//
// # Security Considerations
//
//   - Use 32-byte keys (AES-256) for production environments
//   - Store keys securely using environment variables or secrets managers
//   - Never hardcode encryption keys in source code
//   - The ChaCha8 PRNG is suitable for nonce generation in GCM mode
//   - Always use HTTPS to prevent token interception
//   - Consider implementing token expiration for additional security
//
// # Token Format
//
// The AEADEncryptor produces tokens with the following structure:
//
//	base64(nonce || ciphertext || authentication_tag)
//
// The nonce is prepended to the ciphertext, and the authentication tag is
// appended by the GCM mode. This ensures both confidentiality and authenticity.
package encryption
