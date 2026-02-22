package gorm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGorm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gorm Suite")
}
