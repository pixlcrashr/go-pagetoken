package pagetoken

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/pixlcrashr/go-pagetoken/encryption"
)

type Checksumer interface {
	Checksum() (uint32, error)
}

type Encodable interface {
	Params() []string
}

type CursorField struct {
	Path  string
	Order Order
	Value string
}

type Cursor struct {
	checksum  uint32
	encryptor encryption.Encryptor
	fields    []CursorField
}

func (b *Cursor) Checksum() uint32 {
	return b.checksum
}

func (t *Cursor) tokenize(d []string) (string, error) {
	bs := bytes.NewBuffer(nil)
	if err := json.NewEncoder(bs).Encode(d); err != nil {
		return "", err
	}

	return t.encryptor.Encrypt(bs.Bytes())
}

var ErrFieldNotFound = errors.New("field not found")

func (c *Cursor) Field(k string) (CursorField, error) {
	for _, param := range c.fields {
		if param.Path == k {
			return param, nil
		}
	}

	return CursorField{}, ErrFieldNotFound
}

func (c *Cursor) Fields() []CursorField {
	return c.fields
}

func (c *Cursor) Next(opts ...CursorOpt) *Cursor {
	newC := &Cursor{
		fields: []CursorField{},
	}

	newC.encryptor = c.encryptor
	newC.checksum = c.checksum

	for _, opt := range opts {
		opt(newC)
	}

	return newC
}

func (c *Cursor) String() (string, error) {
	d := make([]string, len(c.fields)*3)

	for i, field := range c.fields {
		d[i*3] = field.Path
		d[i*3+1] = field.Value
		d[i*3+2] = field.Order.String()
	}

	return c.tokenize(append(d, strconv.FormatUint(uint64(c.checksum), 10)))
}

type CursorOpt func(*Cursor)

// CursorField adds a field to the cursor.
//
// The field will be included in the cursor token.
func Field(key, value string, order Order) CursorOpt {
	return func(c *Cursor) {
		c.fields = append(c.fields, CursorField{
			Path:  key,
			Value: value,
			Order: order,
		})
	}
}
