package pagetoken_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pixlcrashr/go-pagetoken"
)

var _ = Describe("Order", func() {
	It("should parse asc", func() {
		Expect(pagetoken.ParseOrder("1")).To(Equal(pagetoken.OrderAsc))
	})

	It("should parse desc", func() {
		Expect(pagetoken.ParseOrder("0")).To(Equal(pagetoken.OrderDesc))
	})

	It("should fail to parse invalid", func() {
		_, err := pagetoken.ParseOrder("invalid")
		Expect(err).To(HaveOccurred())
	})

	It("should return asc string", func() {
		Expect(pagetoken.OrderAsc.String()).To(Equal("1"))
	})

	It("should return desc string", func() {
		Expect(pagetoken.OrderDesc.String()).To(Equal("0"))
	})
})
