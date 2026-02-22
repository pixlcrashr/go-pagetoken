// Package checksum provides CRC32-based checksum validation for page tokens.
//
// This package implements checksum generation and validation to detect when
// pagination parameters change between requests. It uses CRC32 with a configurable
// mask to create checksums from field data, preventing clients from changing
// filters, sorts, or other parameters mid-pagination.
//
// # Features
//
//   - CRC32 checksums with configurable mask
//   - Builder pattern for flexible checksum construction
//   - Field-based checksum generation
//   - Default mask for common use cases
//
// # Purpose
//
// Checksums prevent pagination inconsistencies. For example, if a client requests
// page 1 with status=active, they cannot use the returned page token with
// status=inactive for page 2. The checksum mismatch will be detected and the
// request rejected.
//
// # Example: Basic Checksum Generation
//
//	package main
//
//	import (
//	    "fmt"
//	    "log"
//
//	    "github.com/pixlcrashr/go-pagetoken/checksum"
//	)
//
//	func main() {
//	    // Create a checksum from multiple fields
//	    builder := checksum.NewBuilder(
//	        checksum.Field("status", "active"),
//	        checksum.Field("category", "electronics"),
//	        checksum.Field("limit", "10"),
//	    )
//
//	    crc, err := builder.Build()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Checksum: 0x%08x\n", crc)
//	}
//
// # Example: Custom Mask
//
//	package main
//
//	import (
//	    "fmt"
//	    "log"
//
//	    "github.com/pixlcrashr/go-pagetoken/checksum"
//	)
//
//	func main() {
//	    // Use a custom mask instead of the default
//	    builder := checksum.NewBuilder(
//	        checksum.Mask(0x12345678),
//	        checksum.Field("user_id", "user123"),
//	        checksum.Field("role", "admin"),
//	    )
//
//	    crc, err := builder.Build()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Checksum with custom mask: 0x%08x\n", crc)
//	}
//
// # Example: Checksum Validation
//
//	package main
//
//	import (
//	    "fmt"
//	    "log"
//
//	    "github.com/pixlcrashr/go-pagetoken/checksum"
//	)
//
//	func main() {
//	    // Original request checksum
//	    builder1 := checksum.NewBuilder(
//	        checksum.Field("status", "active"),
//	        checksum.Field("limit", "10"),
//	    )
//	    originalCRC, _ := builder1.Build()
//
//	    // New request with same parameters
//	    builder2 := checksum.NewBuilder(
//	        checksum.Field("status", "active"),
//	        checksum.Field("limit", "10"),
//	    )
//	    newCRC, _ := builder2.Build()
//
//	    // Validate checksums match
//	    if originalCRC == newCRC {
//	        fmt.Println("Checksums match - parameters unchanged")
//	    } else {
//	        fmt.Println("Checksums differ - parameters changed!")
//	    }
//
//	    // Request with different parameters
//	    builder3 := checksum.NewBuilder(
//	        checksum.Field("status", "inactive"), // Changed!
//	        checksum.Field("limit", "10"),
//	    )
//	    changedCRC, _ := builder3.Build()
//
//	    if originalCRC != changedCRC {
//	        log.Fatal("Parameter mismatch detected")
//	    }
//	}
//
// # Example: Integration with Request Types
//
//	package main
//
//	import (
//	    "fmt"
//
//	    "github.com/pixlcrashr/go-pagetoken/checksum"
//	)
//
//	type SearchRequest struct {
//	    Query     string
//	    Category  string
//	    MinPrice  int
//	    MaxPrice  int
//	    SortBy    string
//	    PageToken string
//	}
//
//	func (r *SearchRequest) GetChecksumFields() []checksum.BuilderOpt {
//	    return []checksum.BuilderOpt{
//	        checksum.Field("query", r.Query),
//	        checksum.Field("category", r.Category),
//	        checksum.Field("min_price", fmt.Sprintf("%d", r.MinPrice)),
//	        checksum.Field("max_price", fmt.Sprintf("%d", r.MaxPrice)),
//	        checksum.Field("sort_by", r.SortBy),
//	    }
//	}
//
//	func main() {
//	    req := &SearchRequest{
//	        Query:    "laptop",
//	        Category: "electronics",
//	        MinPrice: 500,
//	        MaxPrice: 2000,
//	        SortBy:   "price_asc",
//	    }
//
//	    // Build checksum from request
//	    builder := checksum.NewBuilder(req.GetChecksumFields()...)
//	    crc, _ := builder.Build()
//
//	    fmt.Printf("Request checksum: 0x%08x\n", crc)
//	}
//
// # Example: Field Order Matters
//
//	package main
//
//	import (
//	    "fmt"
//
//	    "github.com/pixlcrashr/go-pagetoken/checksum"
//	)
//
//	func main() {
//	    // Order 1
//	    builder1 := checksum.NewBuilder(
//	        checksum.Field("a", "1"),
//	        checksum.Field("b", "2"),
//	    )
//	    crc1, _ := builder1.Build()
//
//	    // Order 2 - same fields, different order
//	    builder2 := checksum.NewBuilder(
//	        checksum.Field("b", "2"),
//	        checksum.Field("a", "1"),
//	    )
//	    crc2, _ := builder2.Build()
//
//	    // Checksums will be different!
//	    fmt.Printf("CRC1: 0x%08x\n", crc1)
//	    fmt.Printf("CRC2: 0x%08x\n", crc2)
//	    fmt.Printf("Equal: %v\n", crc1 == crc2) // false
//
//	    // Always use consistent field ordering in your application
//	}
//
// # Default Mask
//
// The default checksum mask is 0x58AEF322. This mask is XORed with the CRC32
// result to provide additional entropy and prevent checksum collisions with
// simple inputs.
//
// # Best Practices
//
//   - Include all parameters that affect query results in the checksum
//   - Use consistent field ordering across all requests
//   - Do not include the page token itself in the checksum
//   - Consider using a custom mask unique to your application
//   - Field names and values are case-sensitive
//
// # Algorithm
//
// The checksum is computed as:
//
//	checksum = CRC32(JSON.encode([field1_key, field1_value, ...])) XOR mask
//
// Where the fields are JSON-encoded as an array of alternating keys and values,
// then hashed with CRC32, and finally XORed with the mask value.
package checksum
