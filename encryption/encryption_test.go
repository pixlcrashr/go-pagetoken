package encryption_test

import (
	"fmt"
	"math/rand/v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pixlcrashr/go-pagetoken/encryption"
)

func randKey(size int) []byte {
	r := rand.ChaCha8{}
	b := make([]byte, size)
	_, err := r.Read(b)
	Expect(err).ToNot(HaveOccurred())
	return b
}

var keySizes = []int{16, 24, 32}

var _ = Describe("Encryption", func() {

	for _, keySize := range keySizes {
		keySize := keySize

		Describe(fmt.Sprintf("with key size %d", keySize), func() {
			var err error
			var e encryption.Encryptor

			BeforeEach(func() {
				e, err = encryption.NewAEADEncryptor(randKey(keySize))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should encrypt a value", func() {
				_, err = e.Encrypt([]byte{1})
				Expect(err).ToNot(HaveOccurred())
			})

			It("should encrypt from and decrypt to the same value", func() {
				in := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
				d, err := e.Encrypt(in)
				Expect(err).ToNot(HaveOccurred())
				out, err := e.Decrypt(d)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(Equal(in))
			})
		})
	}
})
