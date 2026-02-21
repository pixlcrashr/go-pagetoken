package checksum_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pixlcrashr/go-pagetoken/checksum"
)

var _ = Describe("Checksum", func() {
	It("should produce the same checksum for the same fields", func() {
		cb1 := checksum.NewBuilder()
		checksum.Field("key1", "value1")(cb1)
		checksum.Field("key2", "value2")(cb1)
		crc1, err := cb1.Build()
		Expect(err).ToNot(HaveOccurred())

		cb2 := checksum.NewBuilder()
		checksum.Field("key1", "value1")(cb2)
		checksum.Field("key2", "value2")(cb2)
		crc2, err := cb2.Build()
		Expect(err).ToNot(HaveOccurred())

		Expect(crc1).To(Equal(crc2))
	})

	It("should not equal if different mask is used", func() {
		cb1 := checksum.NewBuilder()
		checksum.Field("key1", "value1")(cb1)
		checksum.Field("key2", "value2")(cb1)
		crc1, err := cb1.Build()
		Expect(err).ToNot(HaveOccurred())

		cb2 := checksum.NewBuilder()
		checksum.Mask(0x12345678)(cb2)
		checksum.Field("key1", "value1")(cb2)
		checksum.Field("key2", "value2")(cb2)
		crc2, err := cb2.Build()
		Expect(err).ToNot(HaveOccurred())

		Expect(crc1).ToNot(Equal(crc2))
	})

	It("should not equal if different value order is used", func() {
		cb1 := checksum.NewBuilder()
		checksum.Field("key1", "value1")(cb1)
		checksum.Field("key2", "value2")(cb1)
		crc1, err := cb1.Build()
		Expect(err).ToNot(HaveOccurred())

		cb2 := checksum.NewBuilder()
		checksum.Field("key2", "value2")(cb2)
		checksum.Field("key1", "value1")(cb2)
		crc2, err := cb2.Build()
		Expect(err).ToNot(HaveOccurred())

		Expect(crc1).ToNot(Equal(crc2))
	})

	It("should not equal if different values are used", func() {
		cb1 := checksum.NewBuilder()
		checksum.Field("key1", "value1")(cb1)
		checksum.Field("key2", "value2")(cb1)
		crc1, err := cb1.Build()
		Expect(err).ToNot(HaveOccurred())

		cb2 := checksum.NewBuilder()
		checksum.Field("key1", "value1")(cb2)
		checksum.Field("key2", "differentvalue!")(cb2)
		crc2, err := cb2.Build()
		Expect(err).ToNot(HaveOccurred())

		Expect(crc1).ToNot(Equal(crc2))
	})

	It("should equal if default", func() {
		cb1 := checksum.NewBuilder()
		crc1, err := cb1.Build()
		Expect(err).ToNot(HaveOccurred())

		cb2 := checksum.NewBuilder()
		crc2, err := cb2.Build()
		Expect(err).ToNot(HaveOccurred())

		Expect(crc1).To(Equal(crc2))
	})
})
