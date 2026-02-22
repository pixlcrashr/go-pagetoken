package pagetoken

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/pixlcrashr/go-pagetoken/encryption"
)

type KeysetField struct {
	Path  string
	Order Order
	Value string
}

type KeysetToken struct {
	checksum uint32
	e        encryption.Crypter
	fields   []KeysetField
}

func (b *KeysetToken) Checksum() uint32 {
	return b.checksum
}

func (t *KeysetToken) tokenize(d []string) (string, error) {
	bs := bytes.NewBuffer(nil)
	if err := json.NewEncoder(bs).Encode(d); err != nil {
		return "", err
	}

	return t.e.Encrypt(bs.Bytes())
}

var ErrFieldNotFound = errors.New("field not found")

func (c *KeysetToken) Field(k string) (KeysetField, error) {
	for _, param := range c.fields {
		if param.Path == k {
			return param, nil
		}
	}

	return KeysetField{}, ErrFieldNotFound
}

func (c *KeysetToken) Fields() []KeysetField {
	return c.fields
}

func (c *KeysetToken) Next(opts ...KeysetTokenOpt) *KeysetToken {
	newC := &KeysetToken{
		fields: []KeysetField{},
	}

	newC.e = c.e
	newC.checksum = c.checksum

	for _, opt := range opts {
		opt(newC)
	}

	return newC
}

func (c *KeysetToken) String() (string, error) {
	d := make([]string, len(c.fields)*3)

	for i, field := range c.fields {
		d[i*3] = field.Path
		d[i*3+1] = field.Value
		d[i*3+2] = field.Order.String()
	}

	return c.tokenize(append(d, strconv.FormatUint(uint64(c.checksum), 10)))
}

type KeysetTokenOpt func(*KeysetToken)

// WithKeysetField adds a field to the keyset token.
//
// The field will be included in the keyset token.
func WithKeysetField(key, value string, order Order) KeysetTokenOpt {
	return func(c *KeysetToken) {
		c.fields = append(c.fields, KeysetField{
			Path:  key,
			Value: value,
			Order: order,
		})
	}
}

type KeysetTokenParser struct {
	e encryption.Crypter
}

type KeysetTokenParserOpt func(*KeysetTokenParser)

func WithKeysetTokenEncryptor(e encryption.Crypter) KeysetTokenParserOpt {
	return func(p *KeysetTokenParser) {
		p.e = e
	}
}

func (p *KeysetTokenParser) Parse(token string) (*KeysetToken, error) {
	d, err := p.e.Decrypt(token)
	if err != nil {
		return nil, err
	}

	var ps []string
	if err := json.Unmarshal(d, &ps); err != nil {
		return nil, err
	}

	crc, err := strconv.ParseUint(ps[len(ps)-1], 10, 64)
	if err != nil {
		return nil, err
	}

	fs := []KeysetField{}
	for i := 0; i < len(ps)-1; i += 3 {
		o, err := ParseOrder(ps[i+2])
		if err != nil {
			return nil, err
		}

		fs = append(fs, KeysetField{
			Path:  ps[i],
			Value: ps[i+1],
			Order: o,
		})
	}

	return &KeysetToken{
		checksum: uint32(crc),
		e:        p.e,
		fields:   fs,
	}, nil
}

func NewKeysetTokenParser(opts ...KeysetTokenParserOpt) *KeysetTokenParser {
	p := &KeysetTokenParser{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}
