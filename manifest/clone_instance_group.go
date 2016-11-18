package manifest

import (
	"flag"
	"fmt"

	"github.com/enaml-ops/enaml"
)

// Cloner is a transformation that clones an instance group.
type Cloner struct {
	InstanceGroup string // IG to clone
	Clone         string // name for the copy
}

func (c *Cloner) Apply(dm *enaml.DeploymentManifest) error {
	ig := dm.GetInstanceGroupByName(c.InstanceGroup)
	if ig == nil {
		return fmt.Errorf("couldn't find instance group %s", c.InstanceGroup)
	}

	clone := *ig
	clone.Name = c.Clone
	return dm.AddInstanceGroup(&clone)
}

func (c *Cloner) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("clone", flag.ExitOnError)
	fs.StringVar(&c.InstanceGroup, "instance-group", "", "name of the instance group to clone")
	fs.StringVar(&c.Clone, "clone", "", "the name to use for the copy")
	return fs
}

// CloneTransformation is a TransformationBuilder that builds the
// 'clone' transformation.
func CloneTransformation(args []string) (Transformation, error) {
	c := &Cloner{}
	fs := c.flagSet()
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if c.InstanceGroup == "" {
		return nil, fmt.Errorf("missing required flag -instance-group")
	}
	if c.Clone == "" {
		return nil, fmt.Errorf("missing required flag -clone")
	}

	return c, nil
}
