package main

import (
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
