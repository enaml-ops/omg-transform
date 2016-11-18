package manifest

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
	Network       string
	StaticIPs     []string
	ipsFlag       string
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

	if len(n.StaticIPs) > 0 {
		ig.Networks[0].StaticIPs = n.StaticIPs
	}

	return nil
}

func (n *NetworkMover) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("change-network", flag.ContinueOnError)
	fs.StringVar(&n.InstanceGroup, "instance-group", "", "name of the instance group")
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

	if n.InstanceGroup == "" {
		return nil, errors.New("missing required flag -instance-group")
	}
	if n.Network == "" {
		return nil, errors.New("missing required flag network")
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
