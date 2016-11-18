package manifest

import "github.com/enaml-ops/enaml"

// Transformation is an action applied to a manifest.
type Transformation interface {
	Apply(*enaml.DeploymentManifest) error
}

// TransformationBuilder is a function that builds a transformation from
// a CLI context.
type TransformationBuilder func(args []string) (Transformation, error)
