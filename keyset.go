package pagetoken

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/pixlcrashr/go-pagetoken/encryption"
)

type KeysetValue struct {
	Path  string
	Order Order
	Value string
}

type KeysetToken struct {
	checksum uint32
	e        encryption.Crypter
	payload  *KeysetPayload
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

func (c *KeysetToken) Payload() *KeysetPayload {
	return c.payload
}

func (c *KeysetToken) Next(opts ...KeysetTokenOpt) *KeysetToken {
	newC := &KeysetToken{
		payload: c.payload,
	}

	newC.e = c.e
	newC.checksum = c.checksum

	for _, opt := range opts {
		opt(newC)
	}

	return newC
}

func (c *KeysetToken) String() (string, error) {
	d := make([]string, len(c.payload.vs)*3)

	for i, field := range c.payload.vs {
		d[i*3] = field.Path
		d[i*3+1] = field.Value
		d[i*3+2] = field.Order.String()
	}

	return c.tokenize(append(d, strconv.FormatUint(uint64(c.checksum), 10)))
}

type KeysetTokenOpt func(*KeysetToken)

func WithKeysetPayload(payload *KeysetPayload) KeysetTokenOpt {
	return func(c *KeysetToken) {
		c.payload = payload
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

	vs := []KeysetValue{}
	for i := 0; i < len(ps)-1; i += 3 {
		o, err := ParseOrder(ps[i+2])
		if err != nil {
			return nil, err
		}

		vs = append(vs, KeysetValue{
			Path:  ps[i],
			Value: ps[i+1],
			Order: o,
		})
	}

	return &KeysetToken{
		checksum: uint32(crc),
		e:        p.e,
		payload:  &KeysetPayload{vs: vs},
	}, nil
}

func NewKeysetTokenParser(opts ...KeysetTokenParserOpt) *KeysetTokenParser {
	p := &KeysetTokenParser{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}
