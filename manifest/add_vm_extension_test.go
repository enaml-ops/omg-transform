package manifest

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/enaml-ops/enaml"
)

var _ = Describe("add vm extension", func() {

	Context("when creating the transformation", func() {
		It("returns an error if no arguments are provided", func() {
			_, err := AddVMExtensionTransformation(nil)
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the instance-group argument is missing", func() {
			_, err := AddVMExtensionTransformation([]string{"-name", "public-lbs"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the extension name is invalid", func() {
			_, err := AddVMExtensionTransformation([]string{"-instance-group", "foo", "-name", ",,"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the name argument is missing", func() {
			_, err := AddVMExtensionTransformation([]string{"-instance-group", "redis-master"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns a transformation when given valid args", func() {
			t, err := AddVMExtensionTransformation([]string{"-instance-group", "foo", "-name", "public-lbs1,public-lbs2"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())
		})

		It("returns an error when given invalid args", func() {
			_, err := AddVMExtensionTransformation([]string{"-instance-group"})
			Ω(err).Should(HaveOccurred())
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

		It("no existing instance group should result in error", func() {
			ve := VMExtension{
				InstanceGroup: "blahblah",
				Name:          "public-lbs",
			}
			Ω(ve.Apply(manifest)).ShouldNot(Succeed())
		})

		It("add vm-extensions to an existing instance group", func() {
			ve := VMExtension{
				InstanceGroup: "mysql_proxy",
				Extensions:    []string{"public-lbs1"},
			}
			Ω(ve.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName("mysql_proxy")
			Ω(ig.VMExtensions).Should(HaveLen(1))
			Ω(ig.VMExtensions[0]).Should(Equal(ve.Extensions[0]))
		})

		It("add vm-extensions to an existing instance group where one extension name with comma", func() {
			ve := VMExtension{
				InstanceGroup: "mysql_proxy",
				Extensions:    []string{"public-lbs1,"},
			}
			Ω(ve.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName("mysql_proxy")
			Ω(ig.VMExtensions).Should(HaveLen(1))
			Ω(ig.VMExtensions[0]).Should(Equal(ve.Extensions[0]))
		})

		It("add vm-extensions to an existing instance group where one extension already exists", func() {
			ve := VMExtension{
				InstanceGroup: "nats",
				Extensions:    []string{"public-lbs1,"},
			}
			Ω(ve.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName("nats")
			Ω(ig.VMExtensions).Should(HaveLen(2))
			Ω(ig.VMExtensions[0]).Should(Equal("test"))
			Ω(ig.VMExtensions[1]).Should(Equal(ve.Extensions[0]))

		})

		It("add vm-extensions to an existing instance group with multiple names", func() {
			ve := VMExtension{
				InstanceGroup: "mysql_proxy",
				Extensions:    []string{"public-lbs1", "public-lbs2"},
			}
			Ω(ve.Apply(manifest)).Should(Succeed())

			ig := manifest.GetInstanceGroupByName("mysql_proxy")
			Ω(ig.VMExtensions).Should(HaveLen(2))
			Ω(ig.VMExtensions[0]).Should(Equal(ve.Extensions[0]))
			Ω(ig.VMExtensions[1]).Should(Equal(ve.Extensions[1]))
		})

	})
})
