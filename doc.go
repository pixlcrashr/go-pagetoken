// Package pagetoken provides secure, encrypted cursor-based pagination for REST APIs.
//
// This package implements cursor-based pagination using encrypted page tokens with
// built-in checksum validation to prevent tampering. It's designed for REST API
// servers that need secure pagination with support for multi-field sorting.
//
// # Features
//
//   - Encrypted page tokens using AES-GCM
//   - Checksum validation to prevent parameter changes between requests
//   - Support for multi-field cursors with sort ordering
//   - Interface-based design for custom encryption implementations
//   - Type-safe API for requests and cursors
//
// # Basic Usage
//
// To use this package, you need to:
//
//  1. Create an encryptor with a secure key
//  2. Implement the Request interface for your API requests
//  3. Use FromRequest to parse or create cursors
//  4. Build next page tokens using cursor.Next()
//
// # Example: Simple Pagination
//
//	package main
//
//	import (
//	    "fmt"
//	    "log"
//	    "time"
//
//	    "github.com/pixlcrashr/go-pagetoken"
//	    "github.com/pixlcrashr/go-pagetoken/checksum"
//	    "github.com/pixlcrashr/go-pagetoken/encryption"
//	)
//
//	// Define your request type
//	type ListUsersRequest struct {
//	    PageToken string
//	    Status    string
//	    Limit     int
//	}
//
//	func (r *ListUsersRequest) GetPageToken() string {
//	    return r.PageToken
//	}
//
//	func (r *ListUsersRequest) GetChecksumFields() []checksum.BuilderOpt {
//	    return []checksum.BuilderOpt{
//	        checksum.Field("status", r.Status),
//	        checksum.Field("limit", fmt.Sprintf("%d", r.Limit)),
//	    }
//	}
//
//	type User struct {
//	    ID        string
//	    Name      string
//	    CreatedAt time.Time
//	}
//
//	func main() {
//	    // Create an encryptor with a secure key
//	    key, err := encryption.Rand32ByteKey()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    encryptor, err := encryption.NewAEADEncryptor(key)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Parse request (empty page token for first page)
//	    req := &ListUsersRequest{
//	        PageToken: "",
//	        Status:    "active",
//	        Limit:     10,
//	    }
//
//	    cursor, err := pagetoken.FromRequest(encryptor, req)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Simulate fetching users
//	    users := []User{
//	        {ID: "user1", Name: "Alice", CreatedAt: time.Now()},
//	        {ID: "user2", Name: "Bob", CreatedAt: time.Now()},
//	    }
//
//	    // Create next page token
//	    if len(users) > 0 {
//	        lastUser := users[len(users)-1]
//	        nextCursor := cursor.Next(
//	            pagetoken.Field("id", lastUser.ID, pagetoken.OrderAsc),
//	            pagetoken.Field("created_at", lastUser.CreatedAt.Format(time.RFC3339), pagetoken.OrderDesc),
//	        )
//
//	        nextToken, err := nextCursor.String()
//	        if err != nil {
//	            log.Fatal(err)
//	        }
//
//	        fmt.Printf("Next page token: %s\n", nextToken)
//	    }
//	}
//
// # Example: Using Cursor Fields
//
//	// Extract cursor fields for database queries
//	cursor, err := pagetoken.FromRequest(encryptor, req)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check if this is a continuation request
//	if len(cursor.Fields()) > 0 {
//	    // Get specific field values
//	    idField, err := cursor.Field("id")
//	    if err == nil {
//	        lastID := idField.Value
//	        order := idField.Order
//	        fmt.Printf("Continue from ID: %s (order: %v)\n", lastID, order)
//	    }
//
//	    // Iterate all fields
//	    for _, field := range cursor.Fields() {
//	        fmt.Printf("Field: %s = %s (order: %v)\n", field.Path, field.Value, field.Order)
//	    }
//	}
//
// # Security
//
// Page tokens are encrypted using AES-GCM and include a checksum of the request
// parameters. This prevents:
//
//   - Token tampering (encryption provides authenticity)
//   - Information disclosure (encryption provides confidentiality)
//   - Parameter changes between pages (checksum validation)
//
// Always use HTTPS in production to prevent token interception.
//
// # Checksum Validation
//
// The checksum prevents clients from changing pagination parameters between
// requests. For example, if a client requests page 1 with status=active, they
// cannot use the returned token with status=inactive for page 2. The checksum
// will detect this and reject the request.
package pagetoken
