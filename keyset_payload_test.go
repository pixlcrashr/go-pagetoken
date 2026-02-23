package pagetoken_test

import (
	"errors"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pixlcrashr/go-pagetoken"
)

// build is a test helper that constructs a KeysetPayload via the builder.
func build(fn func(b *pagetoken.KeysetPayloadBuilder)) *pagetoken.KeysetPayload {
	b := &pagetoken.KeysetPayloadBuilder{}
	fn(b)
	return b.Build()
}

var _ = Describe("KeysetPayload", func() {

	Describe("Values", func() {
		It("returns an empty slice for an empty payload", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).Build()
			Expect(p.Values()).To(BeEmpty())
		})

		It("returns values in insertion order", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("a", "first", pagetoken.OrderAsc).
					AddString("b", "second", pagetoken.OrderDesc)
			})
			Expect(p.Values()).To(HaveLen(2))
			Expect(p.Values()[0].Path).To(Equal("a"))
			Expect(p.Values()[1].Path).To(Equal("b"))
		})
	})

	// --- string ---

	Describe("String", func() {
		It("returns the stored string and order", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("s", "hello", pagetoken.OrderAsc)
			})
			v, o, err := p.String("s")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal("hello"))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).Build()
			_, _, err := p.String("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})
	})

	// --- bool ---

	Describe("Bool", func() {
		DescribeTable("round-trips valid booleans",
			func(v bool, order pagetoken.Order) {
				p := build(func(b *pagetoken.KeysetPayloadBuilder) {
					b.AddBool("b", v, order)
				})
				got, gotOrder, err := p.Bool("b")
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(Equal(v))
				Expect(gotOrder).To(Equal(order))
			},
			Entry("true / asc", true, pagetoken.OrderAsc),
			Entry("false / desc", false, pagetoken.OrderDesc),
		)

		It("returns ErrFieldNotFound for a missing key", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).Build()
			_, _, err := p.Bool("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("b", "not-a-bool", pagetoken.OrderAsc)
			})
			_, _, err := p.Bool("b")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- signed integers ---

	Describe("Int", func() {
		It("round-trips via AddInt", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddInt("n", -42, pagetoken.OrderDesc)
			})
			v, o, err := p.Int("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(-42))
			Expect(o).To(Equal(pagetoken.OrderDesc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Int("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-an-int", pagetoken.OrderAsc)
			})
			_, _, err := p.Int("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Int8", func() {
		It("round-trips via AddInt8", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddInt8("n", -8, pagetoken.OrderAsc)
			})
			v, o, err := p.Int8("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int8(-8)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Int8("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "999", pagetoken.OrderAsc)
			})
			_, _, err := p.Int8("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Int16", func() {
		It("round-trips via AddInt16", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddInt16("n", -1000, pagetoken.OrderDesc)
			})
			v, _, err := p.Int16("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int16(-1000)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Int16("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Int16("n")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "99999", pagetoken.OrderAsc)
			})
			_, _, err := p.Int16("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Int32", func() {
		It("round-trips via AddInt32", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddInt32("n", -100000, pagetoken.OrderAsc)
			})
			v, _, err := p.Int32("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int32(-100000)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Int32("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Int32("n")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "9999999999", pagetoken.OrderAsc)
			})
			_, _, err := p.Int32("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Int64", func() {
		It("round-trips via AddInt64", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddInt64("n", -9000000000, pagetoken.OrderDesc)
			})
			v, o, err := p.Int64("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int64(-9000000000)))
			Expect(o).To(Equal(pagetoken.OrderDesc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Int64("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Int64("n")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "99999999999999999999", pagetoken.OrderAsc)
			})
			_, _, err := p.Int64("n")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- unsigned integers ---

	Describe("Uint", func() {
		It("round-trips via AddUint", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddUint("n", 42, pagetoken.OrderAsc)
			})
			v, o, err := p.Uint("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint(42)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Uint("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-uint", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Uint8", func() {
		It("round-trips via AddUint8", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddUint8("n", 255, pagetoken.OrderDesc)
			})
			v, _, err := p.Uint8("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint8(255)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Uint8("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "999", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint8("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Uint16", func() {
		It("round-trips via AddUint16", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddUint16("n", 1000, pagetoken.OrderAsc)
			})
			v, _, err := p.Uint16("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint16(1000)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Uint16("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint16("n")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "99999", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint16("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Uint32", func() {
		It("round-trips via AddUint32", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddUint32("n", 100000, pagetoken.OrderDesc)
			})
			v, _, err := p.Uint32("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint32(100000)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Uint32("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint32("n")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "9999999999", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint32("n")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Uint64", func() {
		It("round-trips via AddUint64", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddUint64("n", 9000000000, pagetoken.OrderAsc)
			})
			v, o, err := p.Uint64("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint64(9000000000)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Uint64("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint64("n")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("n", "99999999999999999999", pagetoken.OrderAsc)
			})
			_, _, err := p.Uint64("n")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- type aliases ---

	Describe("Byte", func() {
		It("is consistent with Uint8", func() {
			const val byte = 0xAB
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddByte("b", val, pagetoken.OrderAsc)
			})
			byteVal, _, err := p.Byte("b")
			Expect(err).NotTo(HaveOccurred())
			uint8Val, _, err := p.Uint8("b")
			Expect(err).NotTo(HaveOccurred())
			Expect(byteVal).To(Equal(byte(uint8Val)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Byte("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("b", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Byte("b")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("b", "999", pagetoken.OrderAsc)
			})
			_, _, err := p.Byte("b")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Rune", func() {
		It("is consistent with Int32", func() {
			const val rune = '✓'
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddRune("r", val, pagetoken.OrderDesc)
			})
			runeVal, _, err := p.Rune("r")
			Expect(err).NotTo(HaveOccurred())
			int32Val, _, err := p.Int32("r")
			Expect(err).NotTo(HaveOccurred())
			Expect(runeVal).To(Equal(rune(int32Val)))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Rune("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("r", "not-a-number", pagetoken.OrderAsc)
			})
			_, _, err := p.Rune("r")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an out-of-range value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("r", "9999999999", pagetoken.OrderAsc)
			})
			_, _, err := p.Rune("r")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- floating point ---

	Describe("Float32", func() {
		It("round-trips via AddFloat32", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddFloat32("f", 1.5, pagetoken.OrderAsc)
			})
			v, o, err := p.Float32("f")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(float32(1.5)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Float32("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("f", "not-a-float", pagetoken.OrderAsc)
			})
			_, _, err := p.Float32("f")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Float64", func() {
		It("round-trips via AddFloat64", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddFloat64("f", 3.141592653589793, pagetoken.OrderDesc)
			})
			v, o, err := p.Float64("f")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(3.141592653589793))
			Expect(o).To(Equal(pagetoken.OrderDesc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Float64("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("f", "not-a-float", pagetoken.OrderAsc)
			})
			_, _, err := p.Float64("f")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- complex ---

	Describe("Complex64", func() {
		It("round-trips via AddComplex64", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddComplex64("c", 1+2i, pagetoken.OrderAsc)
			})
			v, o, err := p.Complex64("c")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(complex64(1 + 2i)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Complex64("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("c", "not-complex", pagetoken.OrderAsc)
			})
			_, _, err := p.Complex64("c")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Complex128", func() {
		It("round-trips via AddComplex128", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddComplex128("c", 3+4i, pagetoken.OrderDesc)
			})
			v, o, err := p.Complex128("c")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(complex128(3 + 4i)))
			Expect(o).To(Equal(pagetoken.OrderDesc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Complex128("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("c", "not-complex", pagetoken.OrderAsc)
			})
			_, _, err := p.Complex128("c")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- time ---

	Describe("Time", func() {
		It("round-trips via AddTime", func() {
			ts := time.Date(2024, 6, 15, 12, 30, 45, 123456789, time.UTC)
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddTime("t", ts, pagetoken.OrderAsc)
			})
			v, o, err := p.Time("t")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(ts))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			_, _, err := (&pagetoken.KeysetPayloadBuilder{}).Build().Time("missing")
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("returns an error for an invalid raw value", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("t", "not-a-time", pagetoken.OrderAsc)
			})
			_, _, err := p.Time("t")
			Expect(err).To(HaveOccurred())
		})
	})

	// --- generic accessor ---

	Describe("GetKeysetValue", func() {
		It("decodes using the supplied function", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("k", "99", pagetoken.OrderAsc)
			})
			v, o, err := pagetoken.GetKeysetValue(p, "k", strconv.Atoi)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(99))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("returns ErrFieldNotFound for a missing key", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).Build()
			_, _, err := pagetoken.GetKeysetValue(p, "missing", strconv.Atoi)
			Expect(err).To(MatchError(pagetoken.ErrFieldNotFound))
		})

		It("propagates a decode error", func() {
			p := build(func(b *pagetoken.KeysetPayloadBuilder) {
				b.AddString("k", "bad", pagetoken.OrderAsc)
			})
			decodeErr := errors.New("bad value")
			_, _, err := pagetoken.GetKeysetValue(p, "k", func(string) (int, error) {
				return 0, decodeErr
			})
			Expect(err).To(MatchError(decodeErr))
		})
	})
})
