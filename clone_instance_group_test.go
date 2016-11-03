package main

import (
	"os"

	"github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("clone instance group", func() {
	Context("PCF 1.8 AWS manifest", func() {
		var manifest *enaml.DeploymentManifest

		BeforeEach(func() {
			f, err := os.Open("fixtures/pcf-aws-1.8.00-build.373.yml")
			Ω(err).ShouldNot(HaveOccurred())
			manifest = enaml.NewDeploymentManifestFromFile(f)
		})

		Context("when cloning a valid instance group", func() {
			const clone = "consul_server_clone"
			BeforeEach(func() {
				c := Cloner{
					InstanceGroup: "consul_server",
					Clone:         clone,
				}
				Ω(c.Apply(manifest)).Should(Succeed())
			})

			It("should have added a new instance group", func() {
				Ω(manifest.GetInstanceGroupByName(clone)).ShouldNot(BeNil())
			})

			It("should have cloned correctly", func() {
				orig := manifest.GetInstanceGroupByName("consul_server")
				clone := manifest.GetInstanceGroupByName(clone)
				Ω(orig).ShouldNot(BeNil())
				Ω(clone).ShouldNot(BeNil())

				By("having the same number of jobs")
				Ω(len(clone.Jobs)).Should(Equal(len(orig.Jobs)))

				By("having identical networks")
				Ω(clone.Networks).Should(Equal(orig.Networks))

				By("having identical lifecycle")
				Ω(clone.Lifecycle).Should(Equal(orig.Lifecycle))

				By("having identical job properties")
				for i := range clone.Jobs {
					if clone.Jobs[i].Properties != nil {
						Ω(clone.Jobs[i].Properties).Should(Equal(orig.Jobs[i].Properties), "job: %s", clone.Jobs[i].Name)
					}
				}
			})
		})
	})
})
