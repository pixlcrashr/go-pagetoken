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

type RequestReader struct {
	e            encryption.Crypter
	checksumOpts []checksum.BuilderOpt
}

type RequestReaderOpt func(*RequestReader)

func WithChecksumOpts(opts ...checksum.BuilderOpt) RequestReaderOpt {
	return func(rr *RequestReader) {
		rr.checksumOpts = opts
	}
}

func WithEncryptor(e encryption.Crypter) RequestReaderOpt {
	return func(rr *RequestReader) {
		rr.e = e
	}
}

func NewRequestReader(
	opts ...RequestReaderOpt,
) *RequestReader {
	rr := &RequestReader{}
	for _, opt := range opts {
		opt(rr)
	}

	// TODO: add defaults
	return rr
}

func (r *RequestReader) createChecksumBuilder(opts ...checksum.BuilderOpt) *checksum.Builder {
	cb := checksum.NewBuilder(opts...)
	return cb
}

func (r *RequestReader) Read(req Request) (*KeysetToken, error) {
	t := req.GetPageToken()
	var c *KeysetToken

	if t == "" {
		// create a newly initialized cursor
		c = &KeysetToken{}
		cb := r.createChecksumBuilder()

		for _, field := range req.GetChecksumFields() {
			field(cb)
		}
		crc, err := cb.Build()
		if err != nil {
			return nil, err
		}
		c.checksum = crc
		c.e = r.e
		c.payload = &KeysetPayload{}
		return c, nil
	}

	p := NewKeysetTokenParser(WithKeysetTokenEncryptor(r.e))
	c, err := p.Parse(t)
	if err != nil {
		return nil, err
	}

	// verify request checksum with page token checksum
	cb := r.createChecksumBuilder()
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
