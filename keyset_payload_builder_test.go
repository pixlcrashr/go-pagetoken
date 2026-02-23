package pagetoken_test

import (
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pixlcrashr/go-pagetoken"
)

var _ = Describe("KeysetPayloadBuilder", func() {

	Describe("Build", func() {
		It("produces an empty payload from a zero-value builder", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).Build()
			Expect(p.Values()).To(BeEmpty())
		})

		It("produces a payload independent of the builder", func() {
			b := &pagetoken.KeysetPayloadBuilder{}
			b.AddString("a", "first", pagetoken.OrderAsc)
			p1 := b.Build()

			b.AddString("b", "second", pagetoken.OrderAsc)
			p2 := b.Build()

			Expect(p1.Values()).To(HaveLen(1), "p1 must not see values added after Build")
			Expect(p2.Values()).To(HaveLen(2))
		})
	})

	Describe("chaining", func() {
		It("preserves insertion order across multiple typed Add calls", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).
				AddString("first", "a", pagetoken.OrderAsc).
				AddInt("second", 2, pagetoken.OrderDesc).
				AddBool("third", true, pagetoken.OrderAsc).
				Build()

			vs := p.Values()
			Expect(vs).To(HaveLen(3))
			Expect(vs[0].Path).To(Equal("first"))
			Expect(vs[1].Path).To(Equal("second"))
			Expect(vs[2].Path).To(Equal("third"))
		})
	})

	// --- string ---

	Describe("AddString", func() {
		It("stores the string value verbatim", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).
				AddString("s", "hello world", pagetoken.OrderAsc).
				Build()
			v, o, err := p.String("s")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal("hello world"))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})
	})

	// --- bool ---

	Describe("AddBool", func() {
		DescribeTable("encodes and round-trips booleans",
			func(val bool, order pagetoken.Order) {
				p := (&pagetoken.KeysetPayloadBuilder{}).AddBool("b", val, order).Build()
				got, gotOrder, err := p.Bool("b")
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(Equal(val))
				Expect(gotOrder).To(Equal(order))
			},
			Entry("true / asc", true, pagetoken.OrderAsc),
			Entry("false / desc", false, pagetoken.OrderDesc),
		)
	})

	// --- signed integers ---

	Describe("AddInt", func() {
		It("round-trips through Int", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddInt("n", -42, pagetoken.OrderDesc).Build()
			v, o, err := p.Int("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(-42))
			Expect(o).To(Equal(pagetoken.OrderDesc))
		})
	})

	Describe("AddInt8", func() {
		It("round-trips through Int8", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddInt8("n", -8, pagetoken.OrderAsc).Build()
			v, _, err := p.Int8("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int8(-8)))
		})
	})

	Describe("AddInt16", func() {
		It("round-trips through Int16", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddInt16("n", -1000, pagetoken.OrderAsc).Build()
			v, _, err := p.Int16("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int16(-1000)))
		})
	})

	Describe("AddInt32", func() {
		It("round-trips through Int32", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddInt32("n", -100000, pagetoken.OrderDesc).Build()
			v, _, err := p.Int32("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int32(-100000)))
		})
	})

	Describe("AddInt64", func() {
		It("round-trips through Int64", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddInt64("n", -9000000000, pagetoken.OrderAsc).Build()
			v, _, err := p.Int64("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(int64(-9000000000)))
		})
	})

	// --- unsigned integers ---

	Describe("AddUint", func() {
		It("round-trips through Uint", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddUint("n", 42, pagetoken.OrderAsc).Build()
			v, o, err := p.Uint("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint(42)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})
	})

	Describe("AddUint8", func() {
		It("round-trips through Uint8", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddUint8("n", 255, pagetoken.OrderDesc).Build()
			v, _, err := p.Uint8("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint8(255)))
		})
	})

	Describe("AddUint16", func() {
		It("round-trips through Uint16", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddUint16("n", 1000, pagetoken.OrderAsc).Build()
			v, _, err := p.Uint16("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint16(1000)))
		})
	})

	Describe("AddUint32", func() {
		It("round-trips through Uint32", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddUint32("n", 100000, pagetoken.OrderDesc).Build()
			v, _, err := p.Uint32("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint32(100000)))
		})
	})

	Describe("AddUint64", func() {
		It("round-trips through Uint64", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddUint64("n", 9000000000, pagetoken.OrderAsc).Build()
			v, _, err := p.Uint64("n")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(uint64(9000000000)))
		})
	})

	// --- type aliases ---

	Describe("AddByte", func() {
		It("is equivalent to AddUint8", func() {
			const val byte = 0xAB
			pByte := (&pagetoken.KeysetPayloadBuilder{}).AddByte("b", val, pagetoken.OrderAsc).Build()
			pUint8 := (&pagetoken.KeysetPayloadBuilder{}).AddUint8("b", val, pagetoken.OrderAsc).Build()

			v1, _, _ := pByte.Byte("b")
			v2, _, _ := pUint8.Uint8("b")
			Expect(v1).To(Equal(byte(v2)))
		})
	})

	Describe("AddRune", func() {
		It("is equivalent to AddInt32", func() {
			const val rune = '✓'
			pRune := (&pagetoken.KeysetPayloadBuilder{}).AddRune("r", val, pagetoken.OrderDesc).Build()
			pInt32 := (&pagetoken.KeysetPayloadBuilder{}).AddInt32("r", val, pagetoken.OrderDesc).Build()

			v1, _, _ := pRune.Rune("r")
			v2, _, _ := pInt32.Int32("r")
			Expect(v1).To(Equal(rune(v2)))
		})
	})

	// --- floating point ---

	Describe("AddFloat32", func() {
		It("round-trips through Float32", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddFloat32("f", 1.5, pagetoken.OrderAsc).Build()
			v, o, err := p.Float32("f")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(float32(1.5)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})
	})

	Describe("AddFloat64", func() {
		It("round-trips through Float64", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddFloat64("f", 3.141592653589793, pagetoken.OrderDesc).Build()
			v, _, err := p.Float64("f")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(3.141592653589793))
		})
	})

	// --- time ---

	Describe("AddTime", func() {
		It("round-trips through Time", func() {
			ts := time.Date(2024, 6, 15, 12, 30, 45, 123456789, time.UTC)
			p := (&pagetoken.KeysetPayloadBuilder{}).AddTime("t", ts, pagetoken.OrderDesc).Build()
			v, o, err := p.Time("t")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(ts))
			Expect(o).To(Equal(pagetoken.OrderDesc))
		})
	})

	// --- complex ---

	Describe("AddComplex64", func() {
		It("round-trips through Complex64", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddComplex64("c", 1+2i, pagetoken.OrderAsc).Build()
			v, o, err := p.Complex64("c")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(complex64(1 + 2i)))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})
	})

	Describe("AddComplex128", func() {
		It("round-trips through Complex128", func() {
			p := (&pagetoken.KeysetPayloadBuilder{}).AddComplex128("c", 3+4i, pagetoken.OrderDesc).Build()
			v, _, err := p.Complex128("c")
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(complex128(3 + 4i)))
		})
	})

	// --- generic adder ---

	Describe("AddKeysetValue", func() {
		It("encodes using the supplied function and round-trips through GetKeysetValue", func() {
			b := &pagetoken.KeysetPayloadBuilder{}
			pagetoken.AddKeysetValue(b, "n", 99, pagetoken.OrderAsc, strconv.Itoa)
			p := b.Build()

			v, o, err := pagetoken.GetKeysetValue(p, "n", strconv.Atoi)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(99))
			Expect(o).To(Equal(pagetoken.OrderAsc))
		})

		It("supports chaining by returning the builder", func() {
			p := pagetoken.AddKeysetValue(
				&pagetoken.KeysetPayloadBuilder{},
				"a", "hello", pagetoken.OrderAsc,
				func(s string) string { return s },
			).AddString("b", "world", pagetoken.OrderDesc).Build()

			Expect(p.Values()).To(HaveLen(2))
		})
	})
})
