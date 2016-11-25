package manifest

import (
	"errors"
	"flag"
	"fmt"

	"github.com/enaml-ops/enaml"
)

//ScaleInstance Scale instance type stores what instance group and how much to scale it
type ScaleInstance struct {
	InstanceGroup string
	Scale         int
}

//Apply apply the scale
func (s *ScaleInstance) Apply(dm *enaml.DeploymentManifest) error {
	ig := dm.GetInstanceGroupByName(s.InstanceGroup)
	if ig == nil {
		return fmt.Errorf("couldn't find instance group %s", s.InstanceGroup)
	}

	ig.Instances = s.Scale

	return nil
}

func (s *ScaleInstance) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("scale", flag.ContinueOnError)
	fs.StringVar(&s.InstanceGroup, "instance-group", "", "name of the instance group")
	fs.IntVar(&s.Scale, "instances", -1, "number of instances")

	return fs
}

//ScaleInstanceTransform Function to scale the instances in the group.
func ScaleInstanceTransform(args []string) (Transformation, error) {
	s := &ScaleInstance{}
	fs := s.flagSet()

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if s.InstanceGroup == "" {
		return nil, errors.New("Missing required flag -instance-group")
	}

	if s.Scale < 0 {
		return nil, errors.New("Missing required flag -scale or invalid value")
	}

	if s.InstanceGroup == "clock_global" {

		if s.Scale > 1 {
			return nil, errors.New("Singleton Instance cannot be scaled higher than 1")
		}
	}

	return s, err
}
