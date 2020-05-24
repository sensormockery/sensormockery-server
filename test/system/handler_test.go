package system_test

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	apiV1 "github.com/sensormockery/sensormockery-server/pkg/api/v1"
)

var _ = Describe("Handler", func() {
	var apiV1URL string

	BeforeEach(func() {
		apiV1URL = fmt.Sprintf("%s%s", os.Getenv(APIURL), apiV1.APIPrefix)
	})

	Describe("Calling to api v1", func() {
		Context("with an invalid path", func() {
			It("should return status not found", func() {
				resp, err := http.Get(apiV1URL + "unsupported")

				if err != nil {
					Fail(err.Error())
				}

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("with an invalid method", func() {
			It("should return status bad request", func() {
				resp, err := http.Get(apiV1URL + "stream")

				if err != nil {
					Fail(err.Error())
				}

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
