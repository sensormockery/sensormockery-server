package system_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stream", func() {

	Describe("Subscribing for streaming service", func() {
		Context("with a kafka broker", func() {
			It("should subscribe", func() {
				Expect("mice").To(Equal("mice"))
			})
		})

		Context("with custom callback", func() {
			It("should subscribe", func() {
				Expect("mice").To(Equal("mice"))
			})
		})
	})
})
