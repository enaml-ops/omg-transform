package main

import (
	"os"

	"github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scale instances", func() {

	Context("When creating the transform", func() {
		It("returns an error if no arguments are provided", func() {
			_, err := ScaleInstanceTransform(nil)
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instance-group argument is missing", func() {
			_, err := ScaleInstanceTransform([]string{"-instances", "1"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instances argument is missing", func() {
			_, err := ScaleInstanceTransform([]string{"-instance-group", "foo"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instances argument is invalid", func() {
			_, err := ScaleInstanceTransform([]string{"-instance-group", "foo", "-instances", "foo"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instances argument is negative", func() {
			_, err := ScaleInstanceTransform([]string{"-instance-group", "foo", "-instances", "-2"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns a transformation when given valid args ", func() {
			t, err := ScaleInstanceTransform([]string{"-instance-group", "foo", "-instances", "1"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())
		})

	})

	Context("PCF 1.8 AWS manifest", func() {
		var (
			manifest *enaml.DeploymentManifest
		)

		BeforeEach(func() {
			f, err := os.Open("fixtures/pcf-aws-1.8.00-build.373.yml")
			Ω(err).ShouldNot(HaveOccurred())
			manifest = enaml.NewDeploymentManifestFromFile(f)

		})

		It("changes the number of instances for a instance", func() {
			const newInstanceCount = 3
			s := ScaleInstance{
				InstanceGroup: "router",
				Scale:         newInstanceCount,
			}
			Ω(s.Apply(manifest)).Should(Succeed())
			ig := manifest.GetInstanceGroupByName("router")
			Ω(ig.Instances).Should(Equal(3))
		})
	})

})
