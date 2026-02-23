package pagetoken

import (
	"strconv"
	"time"
)

type KeysetPayload struct {
	vs []KeysetValue
}

func (kf *KeysetPayload) Values() []KeysetValue {
	return kf.vs
}

// value looks up a single value by path name.
func (kf *KeysetPayload) value(key string) (KeysetValue, error) {
	for _, f := range kf.vs {
		if f.Path == key {
			return f, nil
		}
	}
	return KeysetValue{}, ErrFieldNotFound
}

func identity[T any](v T) (T, error) {
	return v, nil
}

// --- string ---

func (kf *KeysetPayload) String(key string) (string, Order, error) {
	return GetKeysetValue(kf, key, identity[string])
}

// --- bool ---

func (kf *KeysetPayload) Bool(key string) (bool, Order, error) {
	return GetKeysetValue(kf, key, strconv.ParseBool)
}

// --- signed integers ---

func (kf *KeysetPayload) Int(key string) (int, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (int, error) {
		v, err := strconv.ParseInt(s, 10, 0)
		return int(v), err
	})
}

func (kf *KeysetPayload) Int8(key string) (int8, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (int8, error) {
		v, err := strconv.ParseInt(s, 10, 8)
		return int8(v), err
	})
}

func (kf *KeysetPayload) Int16(key string) (int16, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (int16, error) {
		v, err := strconv.ParseInt(s, 10, 16)
		return int16(v), err
	})
}

func (kf *KeysetPayload) Int32(key string) (int32, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (int32, error) {
		v, err := strconv.ParseInt(s, 10, 32)
		return int32(v), err
	})
}

func (kf *KeysetPayload) Int64(key string) (int64, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (int64, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		return v, err
	})
}

// --- unsigned integers ---

func (kf *KeysetPayload) Uint(key string) (uint, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (uint, error) {
		v, err := strconv.ParseUint(s, 10, 0)
		return uint(v), err
	})
}

func (kf *KeysetPayload) Uint8(key string) (uint8, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (uint8, error) {
		v, err := strconv.ParseUint(s, 10, 8)
		return uint8(v), err
	})
}

func (kf *KeysetPayload) Uint16(key string) (uint16, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (uint16, error) {
		v, err := strconv.ParseUint(s, 10, 16)
		return uint16(v), err
	})
}

func (kf *KeysetPayload) Uint32(key string) (uint32, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (uint32, error) {
		v, err := strconv.ParseUint(s, 10, 32)
		return uint32(v), err
	})
}

func (kf *KeysetPayload) Uint64(key string) (uint64, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (uint64, error) {
		v, err := strconv.ParseUint(s, 10, 64)
		return v, err
	})
}

// --- type aliases (byte = uint8, rune = int32) ---

// Byte is an alias accessor for uint8.
func (kf *KeysetPayload) Byte(key string) (byte, Order, error) {
	return kf.Uint8(key)
}

// Rune is an alias accessor for int32.
func (kf *KeysetPayload) Rune(key string) (rune, Order, error) {
	return kf.Int32(key)
}

// --- floating point ---

func (kf *KeysetPayload) Float32(key string) (float32, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (float32, error) {
		v, err := strconv.ParseFloat(s, 32)
		return float32(v), err
	})
}

func (kf *KeysetPayload) Float64(key string) (float64, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (float64, error) {
		v, err := strconv.ParseFloat(s, 64)
		return v, err
	})
}

// --- complex ---

func (kf *KeysetPayload) Complex64(key string) (complex64, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (complex64, error) {
		v, err := strconv.ParseComplex(s, 64)
		return complex64(v), err
	})
}

func (kf *KeysetPayload) Complex128(key string) (complex128, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (complex128, error) {
		v, err := strconv.ParseComplex(s, 128)
		return v, err
	})
}

// --- time ---

func (kf *KeysetPayload) Time(key string) (time.Time, Order, error) {
	return GetKeysetValue(kf, key, func(s string) (time.Time, error) {
		return time.Parse(time.RFC3339Nano, s)
	})
}

// --- generic accessor ---

type KeysetValueDecodeFn[T any] func(string) (T, error)

// GetKeysetValue retrieves a value from kf, converting the raw string
// to the desired type T using the supplied decode function.
//
// Go does not allow methods with additional type parameters on non-generic
// receiver types, so this is a package-level function.
//
// Example:
//
//	id, order, err := pagetoken.GetKeysetValue(payload, "id", uuid.Parse)
func GetKeysetValue[T any](kf *KeysetPayload, key string, decodeFn KeysetValueDecodeFn[T]) (T, Order, error) {
	f, err := kf.value(key)
	if err != nil {
		var zero T
		return zero, OrderDesc, err
	}

	v, err := decodeFn(f.Value)
	if err != nil {
		var zero T
		return zero, f.Order, err
	}
	return v, f.Order, nil
}
