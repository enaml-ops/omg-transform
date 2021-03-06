// omg-transform is a tool for applying transformations to
// bosh deployment manifests.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-transform/manifest"
	yaml "gopkg.in/yaml.v2"
)

// Version is the version of omg-transform.
var Version = "v0.0.0-localcompile"

func init() {
	RegisterTransformationBuilder("change-network", manifest.ChangeNetworkTransformation)
	RegisterTransformationBuilder("clone", manifest.CloneTransformation)
	RegisterTransformationBuilder("change-az", manifest.ChangeAZTransformation)
	RegisterTransformationBuilder("add-tags", manifest.AddTagsTransformation)
	RegisterTransformationBuilder("add-vm-extension", manifest.AddVMExtensionTransformation)
}

func main() {

	if len(os.Args) == 2 && strings.HasSuffix(os.Args[1], "version") {
		fmt.Fprintf(os.Stdout, "Version: %s \n", Version)
		os.Exit(0)
	}

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <transform> [args...]\n", os.Args[0])
		writeTransforms(os.Stderr)
		os.Exit(1)
	}

	name := os.Args[1]
	builder, ok := transformationBuilders[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "Usage: %s <transform> [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "ERROR: unknown transform %q\n", name)
		writeTransforms(os.Stderr)
		os.Exit(1)
	}

	// create the transform based on the arg passed in by the user
	transform, err := builder(os.Args[2:])
	if err == flag.ErrHelp {
		// help message was printed, so just exit
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	// read manifest from stdin
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	manifest := enaml.NewDeploymentManifest(b)
	if manifest == nil {
		fmt.Fprintln(os.Stderr, "ERROR: invalid input manifest")
		os.Exit(1)
	}

	// apply the transformation
	err = transform.Apply(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	// write the transformed manifest back to stdout
	b, err = yaml.Marshal(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(b)
}

func writeTransforms(w io.Writer) {
	fmt.Fprintf(w, "Transforms:\n")
	for t := range transformationBuilders {
		fmt.Fprintf(w, "  %s\n", t)
	}
}

var transformationBuilders map[string]manifest.TransformationBuilder

// RegisterTransformationBuilder registers a transformation builder with the specified name.
func RegisterTransformationBuilder(name string, tb manifest.TransformationBuilder) {
	if transformationBuilders == nil {
		transformationBuilders = make(map[string]manifest.TransformationBuilder)
	}
	if _, ok := transformationBuilders[name]; ok {
		panic(fmt.Errorf("duplicate transformation %q\n\nThis is a development error and should be reported at https://github.com/enaml-ops/omg-transform/issues", name))
	}
	transformationBuilders[name] = tb
}
