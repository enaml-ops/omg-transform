package main

import (
	"os"

	"github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("change AZ transformation", func() {
	Context("when creating the transformation", func() {

		It("returns an error if no arguments are provided", func() {
			_, err := ChangeAZTransformation(nil)
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instance-group argument is missing", func() {
			_, err := ChangeAZTransformation([]string{"-az", "az1,az2"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the az argument is missing", func() {
			_, err := ChangeAZTransformation([]string{"-instance-group", "foo"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the az argument is malformed with space", func() {
			_, err := ChangeAZTransformation([]string{"-instance-group", "foo", "-az", "az1 az2"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the az argument is malformed with empty set for az", func() {
			_, err := ChangeAZTransformation([]string{"-instance-group", "foo", "-az", ",,"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns a transformation when given valid args", func() {
			t, err := ChangeAZTransformation([]string{"-instance-group", "foo", "-az", "az1,az2"})
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

		It("changes the AZ for an existing partition", func() {
			n := AZChanger{
				InstanceGroup: "router",
				AZs:           []string{"az1", "az2"},
			}
			Ω(n.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName(n.InstanceGroup)
			Ω(ig.AZs).Should(HaveLen(len(n.AZs)))
			for i := range n.AZs {
				Ω(ig.AZs[i]).Should(Equal(n.AZs[i]))
			}
		})

		It("returns an error when supplied with a non-existent partition", func() {
			n := NetworkMover{
				InstanceGroup: "this-instance-group-doesnt-exist",
				Network:       "unused",
			}
			Ω(n.Apply(manifest)).ShouldNot(Succeed())
		})
	})
})
