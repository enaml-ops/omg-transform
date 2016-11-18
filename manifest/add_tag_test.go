package manifest

import (
	"os"

	"github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("add tags transformation", func() {
	Context("when creating the transformation", func() {

		It("returns an error if no arguments are provided", func() {
			_, err := AddTagsTransformation(nil)
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the tag doesn't have an =", func() {
			_, err := AddTagsTransformation([]string{"az"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if the tag is missing a value", func() {
			_, err := AddTagsTransformation([]string{"key="})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if one of the tags doesn't have an =", func() {
			_, err := AddTagsTransformation([]string{"tag1=foo", "tag2=bar", "badtag"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error if one of the tags has extra =", func() {
			_, err := AddTagsTransformation([]string{"tag1=foo=bar"})
			Ω(err).Should(HaveOccurred())
		})

		It("returns a transformation when given valid args", func() {
			t, err := AddTagsTransformation([]string{"key1=value1", "key2=value2"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).ShouldNot(BeNil())

			at, ok := t.(*TagAdder)
			Ω(ok).Should(BeTrue())
			Ω(at.Args).Should(HaveLen(2))
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

		It("adds tags", func() {
			t := &TagAdder{
				Args: []string{"key1=value1", "key2=value2"},
			}
			Ω(t.Apply(manifest)).Should(Succeed())
			Ω(manifest.Tags).Should(HaveLen(2))
			Ω(manifest.Tag("key1")).Should(Equal("value1"))
			Ω(manifest.Tag("key2")).Should(Equal("value2"))
		})
	})
})
