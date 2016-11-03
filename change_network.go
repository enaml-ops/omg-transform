package main

import (
	"errors"
	"flag"
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

func (n *NetworkMover) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("change-network", flag.ExitOnError)
	fs.StringVar(&n.InstanceGroup, "instance-group", "", "name of the instance group")
	fs.StringVar(&n.Network, "network", "", "the name of the network to use")
	return fs
}

// ChangeNetworkTransformation is a TransformationBuilder that builds the
// 'change-network' transformation.
func ChangeNetworkTransformation(args []string) (Transformation, error) {
	n := &NetworkMover{}
	fs := n.flagSet()
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if n.InstanceGroup == "" {
		return nil, errors.New("missing required flag -instance-group")
	}
	if n.Network == "" {
		return nil, errors.New("missing required flag network")
	}
	return n, nil
}
