package main

import "github.com/enaml-ops/enaml"

// Transformation is an action applied to a manifest.
type Transformation interface {
	Apply(*enaml.DeploymentManifest) error
}

// TransformationBuilder is a function that builds a transformation from
// a CLI context.
type TransformationBuilder func(args []string) (Transformation, error)

var transformationBuilders map[string]TransformationBuilder

// RegisterTransformationBuilder registers a transformation builder with the specified name.
func RegisterTransformationBuilder(name string, tb TransformationBuilder) {
	if transformationBuilders == nil {
		transformationBuilders = make(map[string]TransformationBuilder)
	}
	transformationBuilders[name] = tb
}
