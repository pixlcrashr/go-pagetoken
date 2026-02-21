# go-pagetoken

[![Test](https://github.com/pixlcrashr/go-pagetoken/actions/workflows/test.yaml/badge.svg)](https://github.com/pixlcrashr/go-pagetoken/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/pixlcrashr/go-pagetoken.svg)](https://pkg.go.dev/github.com/pixlcrashr/go-pagetoken)
[![Go Report Card](https://goreportcard.com/badge/github.com/pixlcrashr/go-pagetoken)](https://goreportcard.com/report/github.com/pixlcrashr/go-pagetoken)

A secure, encrypted cursor-based pagination library for Go REST APIs.

## Overview

`go-pagetoken` provides a robust implementation of cursor-based pagination using encrypted page tokens. It's designed for REST API servers that need secure, tamper-proof pagination with support for multi-field sorting and filtering.

### Key Features

- **🔒 Secure Encryption**: Page tokens are encrypted using AES-GCM to prevent tampering and information disclosure
- **✓ Checksum Validation**: Built-in CRC32 checksums ensure pagination parameters haven't changed between requests
- **📊 Multi-Field Cursors**: Support for pagination cursors with multiple fields and sort orders
- **🔄 Order Support**: Handle both ascending and descending sort orders
- **🛡️ Type-Safe API**: Strongly-typed interfaces for requests and cursors
- **⚡ Zero Dependencies**: Only standard library and testing dependencies

## Installation

```bash
go get github.com/pixlcrashr/go-pagetoken
```

## Quick Start

### 1. Create an Encryptor

First, create an encryptor with a secure key (16, 24, or 32 bytes for AES-128, AES-192, or AES-256):

```go
import (
    "github.com/pixlcrashr/go-pagetoken"
    "github.com/pixlcrashr/go-pagetoken/encryption"
)

// Option 1: Use a secure key from environment or secrets manager
key := []byte("your-32-byte-secret-key-here!!!!")
encryptor, err := encryption.NewAEADEncryptor(key)
if err != nil {
    log.Fatal(err)
}

// Option 2: Generate a random key (for development/testing)
key, err := encryption.Rand32ByteKey()
if err != nil {
    log.Fatal(err)
}
encryptor, err := encryption.NewAEADEncryptor(key)
if err != nil {
    log.Fatal(err)
}
```

### 2. Implement the Request Interface

Your API request struct must implement the `pagetoken.Request` interface:

```go
import "github.com/pixlcrashr/go-pagetoken/checksum"

type ListUsersRequest struct {
    PageToken string
    Status    string
    Limit     int
}

func (r *ListUsersRequest) GetPageToken() string {
    return r.PageToken
}

func (r *ListUsersRequest) GetChecksumFields() []checksum.BuilderOpt {
    // Include all fields that affect pagination results
    return []checksum.BuilderOpt{
        checksum.Field("status", r.Status),
        checksum.Field("limit", fmt.Sprintf("%d", r.Limit)),
    }
}
```

### 3. Parse and Use the Cursor

```go
func ListUsers(req *ListUsersRequest) ([]User, string, error) {
    // Parse the page token from the request
    cursor, err := pagetoken.FromRequest(encryptor, req)
    if err != nil {
        return nil, "", err
    }

    // Extract cursor fields for your query
    var lastID string
    var lastCreatedAt time.Time

    if len(cursor.Fields()) > 0 {
        idField, _ := cursor.Field("id")
        lastID = idField.Value

        createdAtField, _ := cursor.Field("created_at")
        lastCreatedAt, _ = time.Parse(time.RFC3339, createdAtField.Value)
    }

    // Execute your database query using cursor values
    users := queryUsers(lastID, lastCreatedAt, req.Limit)

    // Create next page token
    if len(users) > 0 {
        lastUser := users[len(users)-1]
        nextCursor := cursor.Next(
            pagetoken.Field("id", lastUser.ID, pagetoken.OrderAsc),
            pagetoken.Field("created_at", lastUser.CreatedAt.Format(time.RFC3339), pagetoken.OrderDesc),
        )

        nextToken, err := nextCursor.String()
        if err != nil {
            return nil, "", err
        }

        return users, nextToken, nil
    }

    return users, "", nil
}
```

### 4. Return the Token to Clients

```go
type ListUsersResponse struct {
    Users         []User `json:"users"`
    NextPageToken string `json:"next_page_token,omitempty"`
}

func HandleListUsers(w http.ResponseWriter, r *http.Request) {
    req := &ListUsersRequest{
        PageToken: r.URL.Query().Get("page_token"),
        Status:    r.URL.Query().Get("status"),
        Limit:     10,
    }

    users, nextToken, err := ListUsers(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(ListUsersResponse{
        Users:         users,
        NextPageToken: nextToken,
    })
}
```

## How It Works

### Encryption Submodule

The library provides a flexible encryption system through the `Encryptor` interface:

- **Interface-Based Design**: The `Encryptor` interface allows for custom encryption implementations
- **AEAD Encryption**: Default implementation uses AES-GCM (Authenticated Encryption with Associated Data)
- **ChaCha8 PRNG**: Nonce generation uses Go's ChaCha8 pseudo-random number generator from `math/rand/v2`
- **Key Generation Helpers**: Built-in functions to generate random keys for different AES variants

**Encryption Flow:**
```
plaintext → JSON encode → AES-GCM encrypt → base64 URL encode → token string
```

**Decryption Flow:**
```
token string → base64 URL decode → AES-GCM decrypt → JSON decode → plaintext
```

### Security Model

1. **Encryption**: Page tokens are encrypted using AES-GCM AEAD, ensuring confidentiality and authenticity
2. **Checksum Validation**: Each token includes a CRC32 checksum of the pagination parameters, preventing clients from changing filters mid-pagination
3. **Tamper Detection**: Any modification to the encrypted token will cause decryption to fail
4. **Nonce Uniqueness**: Each encryption operation uses a fresh nonce generated via ChaCha8 PRNG

### Token Format

Internally, tokens contain:
- Multiple cursor fields (path, value, sort order)
- A CRC32 checksum of the request parameters
- Everything is JSON-encoded, encrypted, and base64-encoded

### Checksum Purpose

The checksum prevents pagination inconsistencies. For example, if a client requests page 1 with `status=active`, they cannot use the returned token with `status=inactive` for page 2. The checksum mismatch will be detected and rejected.

## API Reference

### Core Types

- **`Cursor`**: Represents a pagination cursor with encrypted fields
- **`Parser`**: Parses encrypted page tokens back into cursors
- **`Request`**: Interface that API requests must implement
- **`CursorField`**: A single field in a cursor (path, value, order)
- **`Order`**: Sort order (ascending or descending)

### Encryption Package

- **`Encryptor`**: Interface for encryption/decryption implementations
- **`AEADEncryptor`**: AES-GCM AEAD implementation of the Encryptor interface
- **`NewAEADEncryptor(key)`**: Create a new AEAD encryptor with AES key (16/24/32 bytes)
- **`Rand16ByteKey()`**: Generate a random 16-byte key for AES-128
- **`Rand24ByteKey()`**: Generate a random 24-byte key for AES-192
- **`Rand32ByteKey()`**: Generate a random 32-byte key for AES-256

### Checksum Package

- **`Builder`**: Builds CRC32 checksums from field data
- **`NewBuilder(opts...)`**: Create a new checksum builder
- **`Field(key, value)`**: Add a field to the checksum
- **`Mask(mask)`**: Set a custom checksum mask (default: 0x58AEF322)

### Main Functions

- **`FromRequest(encryptor, request, opts...)`**: Parse or create a cursor from an API request
- **`Field(key, value, order)`**: Create a cursor field option
- **`NewParser(opts...)`**: Create a new token parser

See the [GoDoc](https://pkg.go.dev/github.com/pixlcrashr/go-pagetoken) for complete API documentation.

## Advanced Usage

### Custom Encryption Implementation

You can implement your own encryption strategy by satisfying the `Encryptor` interface:

```go
type Encryptor interface {
    Encrypt(d []byte) (string, error)
    Decrypt(token string) ([]byte, error)
}
```

Example with a custom encryptor:

```go
type CustomEncryptor struct {
    // Your custom fields
}

func (e *CustomEncryptor) Encrypt(data []byte) (string, error) {
    // Your encryption logic
    return encryptedString, nil
}

func (e *CustomEncryptor) Decrypt(token string) ([]byte, error) {
    // Your decryption logic
    return decryptedData, nil
}

// Use with pagetoken
customEnc := &CustomEncryptor{}
cursor, err := pagetoken.FromRequest(customEnc, req)
```

### Key Generation for Production

For production environments, generate and store keys securely:

```go
// Generate a key once
key, err := encryption.Rand32ByteKey()
if err != nil {
    log.Fatal(err)
}

// Store in environment or secrets manager
// Example: export PAGE_TOKEN_KEY=$(echo $key | base64)

// Load in application
keyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("PAGE_TOKEN_KEY"))
if err != nil {
    log.Fatal(err)
}

encryptor, err := encryption.NewAEADEncryptor(keyBytes)
if err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. **Secure Key Management**: Store encryption keys securely (environment variables, secrets manager)
2. **Key Rotation**: Plan for key rotation by versioning your tokens or supporting multiple keys
3. **Include Relevant Fields**: Add all fields that affect query results to the checksum
4. **Consistent Ordering**: Use the same field ordering across requests for predictable pagination
5. **Limit Page Size**: Enforce reasonable page size limits to prevent performance issues
6. **Production Keys**: Always use cryptographically secure random keys in production (use the `Rand*ByteKey()` helpers)

## Security Considerations

- **Key Size**: Use 32-byte keys (AES-256) for production environments
- **Key Storage**: Never hardcode encryption keys; use secure key management
- **Key Generation**: The `Rand*ByteKey()` functions use ChaCha8 PRNG which is suitable for key generation but not cryptographic random number generation in security-critical contexts
- **Nonce Generation**: AEADEncryptor uses ChaCha8 for nonce generation, which is appropriate for GCM mode
- **Token Lifetime**: Consider implementing token expiration for additional security
- **HTTPS Only**: Always use HTTPS to prevent token interception
- **Information Disclosure**: While tokens are encrypted, avoid including sensitive data in cursor fields
- **Custom Implementations**: If implementing a custom `Encryptor`, ensure your implementation provides authenticated encryption

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

See [LICENSE](LICENSE) file for details.
