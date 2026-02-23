package pagetoken

import (
	"strconv"
	"time"
)

// KeysetPayloadBuilder accumulates KeysetValues in insertion order.
// Values may only be appended; the builder provides no way to delete,
// edit, or reorder existing entries.
type KeysetPayloadBuilder struct {
	vs []KeysetValue
}

func NewKeysetPayloadBuilder() *KeysetPayloadBuilder {
	return &KeysetPayloadBuilder{}
}

// append is the single internal write path; all typed adders funnel through here.
func (b *KeysetPayloadBuilder) append(key, value string, order Order) *KeysetPayloadBuilder {
	b.vs = append(b.vs, KeysetValue{
		Path:  key,
		Value: value,
		Order: order,
	})
	return b
}

// Build produces an immutable KeysetPayload from the values accumulated
// so far. The returned payload is independent of the builder: subsequent
// adder calls do not affect it.
func (b *KeysetPayloadBuilder) Build() *KeysetPayload {
	vs := make([]KeysetValue, len(b.vs))
	copy(vs, b.vs)
	return &KeysetPayload{vs: vs}
}

// --- string ---

func (b *KeysetPayloadBuilder) AddString(key string, value string, order Order) *KeysetPayloadBuilder {
	return b.append(key, value, order)
}

// --- bool ---

func (b *KeysetPayloadBuilder) AddBool(key string, value bool, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatBool(value), order)
}

// --- signed integers ---

func (b *KeysetPayloadBuilder) AddInt(key string, value int, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatInt(int64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddInt8(key string, value int8, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatInt(int64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddInt16(key string, value int16, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatInt(int64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddInt32(key string, value int32, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatInt(int64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddInt64(key string, value int64, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatInt(value, 10), order)
}

// --- unsigned integers ---

func (b *KeysetPayloadBuilder) AddUint(key string, value uint, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatUint(uint64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddUint8(key string, value uint8, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatUint(uint64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddUint16(key string, value uint16, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatUint(uint64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddUint32(key string, value uint32, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatUint(uint64(value), 10), order)
}

func (b *KeysetPayloadBuilder) AddUint64(key string, value uint64, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatUint(value, 10), order)
}

// --- type aliases (byte = uint8, rune = int32) ---

// AddByte is an alias adder for uint8.
func (b *KeysetPayloadBuilder) AddByte(key string, value byte, order Order) *KeysetPayloadBuilder {
	return b.AddUint8(key, value, order)
}

// AddRune is an alias adder for int32.
func (b *KeysetPayloadBuilder) AddRune(key string, value rune, order Order) *KeysetPayloadBuilder {
	return b.AddInt32(key, value, order)
}

// --- floating point ---

func (b *KeysetPayloadBuilder) AddFloat32(key string, value float32, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatFloat(float64(value), 'g', -1, 32), order)
}

func (b *KeysetPayloadBuilder) AddFloat64(key string, value float64, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatFloat(value, 'g', -1, 64), order)
}

// --- complex ---

func (b *KeysetPayloadBuilder) AddComplex64(key string, value complex64, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatComplex(complex128(value), 'g', -1, 64), order)
}

func (b *KeysetPayloadBuilder) AddComplex128(key string, value complex128, order Order) *KeysetPayloadBuilder {
	return b.append(key, strconv.FormatComplex(value, 'g', -1, 128), order)
}

// --- time ---

func (b *KeysetPayloadBuilder) AddTime(key string, value time.Time, order Order) *KeysetPayloadBuilder {
	return b.append(key, value.Format(time.RFC3339Nano), order)
}

type KeysetValueEncodeFn[T any] func(T) string

// AddKeysetValue is the inverse of GetKeysetValue: it serialises value to a
// string using the supplied format function and appends it to the builder.
//
// Example:
//
//	pagetoken.AddKeysetValue(b, "id", someUUID, pagetoken.OrderAsc, uuid.UUID.String)
func AddKeysetValue[T any](b *KeysetPayloadBuilder, key string, value T, order Order, encodeFn KeysetValueEncodeFn[T]) *KeysetPayloadBuilder {
	return b.append(key, encodeFn(value), order)
}
