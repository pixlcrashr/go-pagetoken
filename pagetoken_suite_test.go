package pagetoken_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPagetoken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pagetoken Suite")
}
