package main

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/enaml-ops/enaml"
)

var _ = Describe("change network", func() {

	Context("when creating the transformation", func() {
		It("returns an error if no arguments are provided", func() {
			_, err := ChangeNetworkTransformation(nil)
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instance-group and lifecycle argument is missing", func() {
			_, err := ChangeNetworkTransformation([]string{"-network", "net"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if both instance-group and lifecycle are present", func() {
			_, err := ChangeNetworkTransformation([]string{"-lifecycle", "life", "-network", "net", "-instance-group", "foo"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the network argument is missing (select by instance-group)", func() {
			_, err := ChangeNetworkTransformation([]string{"-instance-group", "foo"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if network argument is missing (select by lifecycle)", func() {
			_, err := ChangeNetworkTransformation([]string{"-lifecycle", "life"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns a transformation when given valid args (select by lifecycle)", func() {
			t, err := ChangeNetworkTransformation([]string{"-network", "net", "-lifecycle", "life"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())
		})

		It("returns a transformation when given valid args (select by instance-group)", func() {
			t, err := ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())
		})

		It("returns an error when given invalid static IP ranges", func() {
			_, err := ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", ",,"})
			Ω(err).Should(HaveOccurred())

			_, err = ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "1.2.3.4-1.2.3.10-1.2.3.11"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error when given invalid IP addresses", func() {
			_, err := ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "1.2.3.4-1.2.3.X"})
			Ω(err).Should(HaveOccurred())

			_, err = ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "abc,def"})
			Ω(err).Should(HaveOccurred())

			_, err = ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "10.0.0.0-10.0.0.2,foo-bar"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns a transformation when given valid args (no IPs)", func() {
			t, err := ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())
		})

		It("returns a transformation when given valid args (with IPs)", func() {
			t, err := ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "1.2.3.4-1.2.3.10"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())

			t, err = ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "1.2.3.4-1.2.3.10,1.2.3.15-1.2.3.20"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())

			t, err = ChangeNetworkTransformation([]string{"-instance-group", "foo", "-network", "net", "-static-ips", "1.2.3.4,1.2.3.6"})
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

		It("changes the network name for an existing partition (select by instance-group)", func() {
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

		It("changes the network name for an existing partition (select by lifecycle)", func() {
			const newNetwork = "newNetwork"
			n := NetworkMover{
				Lifecycle: "errand",
				Network:   newNetwork,
			}
			Ω(n.Apply(manifest)).Should(Succeed())

			for _, ig := range manifest.InstanceGroups {
				if ig.Lifecycle == n.Lifecycle {
					Ω(ig.Networks).Should(HaveLen(1))
					Ω(ig.Networks[0].Name).Should(Equal(newNetwork), "Instance group "+ig.Name+" failed to change.")
				}
			}
		})

		It("changes the network's static IPs for an existing partition", func() {
			const newNetwork = "newNetwork"
			n := NetworkMover{
				InstanceGroup: "mysql_proxy",
				Network:       newNetwork,
				StaticIPs:     []string{"10.0.0.3-10.0.0.10", "10.0.16.5"},
			}
			Ω(n.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName("mysql_proxy")
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal(newNetwork))

			Ω(ig.Networks[0].StaticIPs).Should(HaveLen(2))
			Ω(ig.Networks[0].StaticIPs).Should(ConsistOf("10.0.0.3-10.0.0.10", "10.0.16.5"))
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
