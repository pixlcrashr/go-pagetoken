package checksum

import (
	"bytes"
	"encoding/json"
	"hash/crc32"
)

const DefaultChecksumMask = 0x58AEF322

func checksum(data []byte, mask uint32) uint32 {
	return crc32.ChecksumIEEE(data) ^ mask
}

type Builder struct {
	mask   uint32
	fields []string
}

func NewBuilder(opts ...BuilderOpt) *Builder {
	cb := &Builder{
		mask:   DefaultChecksumMask,
		fields: make([]string, 0),
	}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

func (b *Builder) Build() (uint32, error) {
	bs := bytes.NewBuffer(nil)
	if err := json.NewEncoder(bs).Encode(b.fields); err != nil {
		return 0, err
	}

	return checksum(bs.Bytes(), b.mask), nil
}

type BuilderOpt func(*Builder)

func Mask(mask uint32) BuilderOpt {
	return func(b *Builder) {
		b.mask = mask
	}
}

func Field(key, value string) BuilderOpt {
	return func(b *Builder) {
		b.fields = append(b.fields, key, value)
	}
}
