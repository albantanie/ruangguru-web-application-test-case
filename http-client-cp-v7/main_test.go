package main_test

import (
	main "a21hc3NpZ25tZW50"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	Describe("HTTP Client GET", func() {
		When("hit API https://api.quotable.io/quotes/random?limit=3 and decode to struct Quotes", func() {
			It("should return struct Quotes with contents of the hit API", func() {
				res, err := main.ClientGet()
				Expect(err).To(BeNil())
				Expect(res).To(HaveLen(3))

				//in experiment how to check JSON response with random order
			})
		})
	})

	Describe("HTTP Client POST", func() {
		When("hit API https://postman-echo.com/post and decode to struct Postman", func() {
			It("should return struct Postman with contents of the hit API", func() {
				res, err := main.ClientPost()
				Expect(err).To(BeNil())
				Expect(res.Url).To(Equal("https://postman-echo.com/post"))
				Expect(res.Data.Name).To(Equal("Dion"))
				Expect(res.Data.Email).To(Equal("dionbe2022@gmail.com"))
			})
		})
	})
})
