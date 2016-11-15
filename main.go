// omg-transform is a tool for applying transformations to
// bosh deployment manifests.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/enaml-ops/enaml"
	yaml "gopkg.in/yaml.v2"
)

// Version is the version of omg-transform.
var Version = "v0.0.1"

func init() {
	RegisterTransformationBuilder("change-network", ChangeNetworkTransformation)
	RegisterTransformationBuilder("clone", CloneTransformation)
	RegisterTransformationBuilder("change-az", ChangeAZTransformation)
}

func main() {
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
