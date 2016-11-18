package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("transformation registry", func() {
	It("panics if multiple transformations share the same name", func() {
		RegisterTransformationBuilder("transform", nil)
		Ω(func() {
			RegisterTransformationBuilder("transform", nil)
		}).Should(Panic())
	})
})
