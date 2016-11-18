package cloudconfig

import "github.com/enaml-ops/enaml"

// Transformation is an action applied to a cloud config.
type Transformation interface {
	Apply(*enaml.CloudConfigManifest) error
}

// TransformationBuilder is a function that builds a transformation from
// a CLI context.
type TransformationBuilder func(args []string) (Transformation, error)
