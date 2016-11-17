package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/enaml-ops/enaml"
)

// NetworkMover is a transformation that changes which network
// an instance group is placed in.
type NetworkMover struct {
	InstanceGroup string
	Lifecycle     string
	Network       string
	StaticIPs     []string
	ipsFlag       string
}

func (n *NetworkMover) Apply(dm *enaml.DeploymentManifest) error {

	if n.InstanceGroup != "" {
		ig := dm.GetInstanceGroupByName(n.InstanceGroup)
		if ig == nil {
			return fmt.Errorf("couldn't find instance group %s", n.InstanceGroup)
		}
		return n.applyToInstanceGroup(ig)
	}

	if n.Lifecycle != "" {
		var err error
		for _, ig := range dm.InstanceGroups {
			if ig.Lifecycle == n.Lifecycle {
				err = n.applyToInstanceGroup(ig)
				if err != nil {
					return fmt.Errorf("error applying transformation to instance group %s: %v\n", ig.Name, err)
				}
			}
		}
		return nil
	}

	return errors.New("transform was not applied by instance group or lifecycle")

}

func (n *NetworkMover) applyToInstanceGroup(ig *enaml.InstanceGroup) error {
	if l := len(ig.Networks); l != 1 {
		return fmt.Errorf("expected 1 network, found %d", l)
	}

	ig.Networks[0].Name = n.Network

	if len(n.StaticIPs) > 0 {
		ig.Networks[0].StaticIPs = n.StaticIPs
	}
	return nil
}

func (n *NetworkMover) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("change-network", flag.ContinueOnError)
	fs.StringVar(&n.InstanceGroup, "instance-group", "", "apply transformation to the instance group with this name")
	fs.StringVar(&n.Lifecycle, "lifecycle", "", "apply transformation to all instance groups with this lifecycle")
	fs.StringVar(&n.Network, "network", "", "the name of the network to use")
	fs.StringVar(&n.ipsFlag, "static-ips", "", "comma-separated list of static IP ranges to set on the network")
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

	igPresent := n.InstanceGroup == ""
	lifecyclePresent := n.Lifecycle == ""

	if igPresent == lifecyclePresent {
		return nil, errors.New("either -lifecycle or -instance-group must be specified, but not both")
	}
	if n.Network == "" {
		return nil, errors.New("missing required flag -network")
	}
	if n.ipsFlag != "" {
		n.StaticIPs = split(n.ipsFlag, ",")
		if len(n.StaticIPs) == 0 {
			return nil, errors.New("invalid -static-ips flag")
		}
		for _, ipRange := range n.StaticIPs {
			c := strings.Count(ipRange, "-")
			if c > 1 {
				return nil, fmt.Errorf("invalid IP range %q", ipRange)
			}
			parts := strings.Split(ipRange, "-")
			for _, ipStr := range parts {
				if ip := net.ParseIP(ipStr); ip == nil {
					return nil, fmt.Errorf("%q is not a valid IP address", ipStr)
				}
			}
		}
	}
	return n, nil
}
