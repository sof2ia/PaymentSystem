package internal

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validation Test", func() {
	It("should validate all fields successfully", func() {
		c := CreateUserRequest{
			Name:  "First User",
			Age:   20,
			Phone: "+5512934523456",
			Email: "name@gmail.com",
			CPF:   "12345678912",
		}
		err := c.Validation()
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("should validate Age unsuccessfully", func() {
		c := CreateUserRequest{
			Name:  "First User",
			Age:   16,
			Phone: "+5512934523456",
			Email: "name@gmail.com",
			CPF:   "12345678912",
		}
		err := c.Validation()
		Expect(err).Should(HaveOccurred())
	})
	It("should validate Phone unsuccessfully", func() {
		c := CreateUserRequest{
			Name:  "First User",
			Age:   20,
			Phone: "5512934523456",
			Email: "name@gmail.com",
			CPF:   "12345678912",
		}
		err := c.Validation()
		Expect(err).Should(HaveOccurred())
	})
	It("should validate Email unsuccessfully", func() {
		c := CreateUserRequest{
			Name:  "First User",
			Age:   20,
			Phone: "+5512934523456",
			Email: "name@gmail",
			CPF:   "12345678912",
		}
		err := c.Validation()
		Expect(err).Should(HaveOccurred())
	})
	It("should validate CPF successfully", func() {
		c := CreateUserRequest{
			Name:  "First User",
			Age:   20,
			Phone: "+5512934523456",
			Email: "name@gmail.com",
			CPF:   "123456789",
		}
		err := c.Validation()
		Expect(err).Should(HaveOccurred())
	})

})
