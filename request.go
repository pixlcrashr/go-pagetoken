package pagetoken

import (
	"fmt"

	"github.com/pixlcrashr/go-pagetoken/checksum"
	"github.com/pixlcrashr/go-pagetoken/encryption"
)

// Request is an interface that must be implemented by all request types that
// use page tokens.
type Request interface {
	// GetChecksumFields returns a list of functions that build the checksum for the request.
	GetChecksumFields() []checksum.BuilderOpt
	// GetPageToken returns the page token from the request.
	GetPageToken() string
}

func FromRequest(e encryption.Encryptor, req Request, checksumOpts ...checksum.BuilderOpt) (*Cursor, error) {
	t := req.GetPageToken()
	var c *Cursor

	if t == "" {
		// create a newly initialized cursor
		c = &Cursor{}
		cb := checksum.NewBuilder(checksumOpts...)
		for _, field := range req.GetChecksumFields() {
			field(cb)
		}
		crc, err := cb.Build()
		if err != nil {
			return nil, err
		}
		c.checksum = crc
		c.encryptor = e
		c.fields = []CursorField{}
		return c, nil
	}

	p := NewParser(WithEncryptor(e))
	c, err := p.Parse(t)
	if err != nil {
		return nil, err
	}

	// verify request checksum with page token checksum
	cb := checksum.NewBuilder(checksumOpts...)
	for _, field := range req.GetChecksumFields() {
		field(cb)
	}
	crc, err := cb.Build()
	if err != nil {
		return nil, err
	}

	if crc != c.checksum {
		return nil, fmt.Errorf(
			"checksum mismatch (got 0x%x but expected 0x%x)", c.checksum, crc,
		)
	}

	return c, nil
}
