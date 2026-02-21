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

// Use a secure key (32 bytes for AES-256)
key := []byte("your-32-byte-secret-key-here!!!!")
encryptor, err := encryption.NewEncryptor(key)
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

### Security Model

1. **Encryption**: Page tokens are encrypted using AES-GCM, ensuring confidentiality and authenticity
2. **Checksum Validation**: Each token includes a checksum of the pagination parameters, preventing clients from changing filters mid-pagination
3. **Tamper Detection**: Any modification to the encrypted token will cause decryption to fail

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

### Main Functions

- **`FromRequest(encryptor, request, opts...)`**: Parse or create a cursor from an API request
- **`Field(key, value, order)`**: Create a cursor field option
- **`NewParser(opts...)`**: Create a new token parser
- **`NewEncryptor(key)`**: Create a new encryptor with AES key

See the [GoDoc](https://pkg.go.dev/github.com/pixlcrashr/go-pagetoken) for complete API documentation.

## Best Practices

1. **Secure Key Management**: Store encryption keys securely (environment variables, secrets manager)
2. **Key Rotation**: Plan for key rotation by versioning your tokens or supporting multiple keys
3. **Include Relevant Fields**: Add all fields that affect query results to the checksum
4. **Consistent Ordering**: Use the same field ordering across requests for predictable pagination
5. **Limit Page Size**: Enforce reasonable page size limits to prevent performance issues

## Security Considerations

- **Key Size**: Use 32-byte keys (AES-256) for production environments
- **Key Storage**: Never hardcode encryption keys; use secure key management
- **Token Lifetime**: Consider implementing token expiration for additional security
- **HTTPS Only**: Always use HTTPS to prevent token interception
- **Information Disclosure**: While tokens are encrypted, avoid including sensitive data in cursor fields

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

See [LICENSE](LICENSE) file for details.
