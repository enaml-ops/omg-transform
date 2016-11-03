package main

import (
	"fmt"

	"github.com/enaml-ops/enaml"
)

// NetworkMover is a transformation that changes which network
// an instance group is placed in.
type NetworkMover struct {
	InstanceGroup string
	Network       string
}

func (n *NetworkMover) Apply(dm *enaml.DeploymentManifest) error {
	ig := dm.GetInstanceGroupByName(n.InstanceGroup)
	if ig == nil {
		return fmt.Errorf("couldn't find instance group %s", n.InstanceGroup)
	}

	if l := len(ig.Networks); l != 1 {
		return fmt.Errorf("expected 1 network, found %d", l)
	}

	ig.Networks[0].Name = n.Network
	return nil
}
