package main

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/enaml-ops/enaml"
)

var _ = Describe("change network", func() {
	Context("PCF 1.8 AWS manifest", func() {
		var (
			manifest *enaml.DeploymentManifest
		)
		BeforeEach(func() {
			f, err := os.Open("fixtures/pcf-aws-1.8.00-build.373.yml")
			Ω(err).ShouldNot(HaveOccurred())
			manifest = enaml.NewDeploymentManifestFromFile(f)
		})

		It("changes the network for an existing partition", func() {
			const newNetwork = "newNetwork"
			n := NetworkMover{
				InstanceGroup: "mysql_proxy",
				Network:       newNetwork,
			}
			Ω(n.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName("mysql_proxy")
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal(newNetwork))
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
