package order_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pixlcrashr/go-pagetoken/order"
)

var _ = Describe("Order", func() {
	It("should parse asc", func() {
		var fs order.Fields

		Expect(fs.UnmarshalString("title")).To(Succeed())
		Expect(fs).To(HaveLen(1))
		Expect(fs[0].Path).To(Equal("title"))
		Expect(fs[0].Order).To(Equal(order.Asc))
	})

	It("should parse desc", func() {
		var fs order.Fields

		Expect(fs.UnmarshalString("title desc")).To(Succeed())
		Expect(fs).To(HaveLen(1))
		Expect(fs[0].Path).To(Equal("title"))
		Expect(fs[0].Order).To(Equal(order.Desc))
	})

	It("should fail to parse invalid", func() {
		var o order.Order
		Expect(o.UnmarshalString("invalid")).ToNot(Succeed())
	})

	It("should return asc string", func() {
		Expect(order.Asc.String()).To(Equal("asc"))
	})

	It("should return desc string", func() {
		Expect(order.Desc.String()).To(Equal("desc"))
	})
})
