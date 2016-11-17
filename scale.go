package main

import (
	"errors"
	"flag"

	"github.com/enaml-ops/enaml"
)

type ScaleInstance struct {
	InstanceGroup string
	Scale         int
}

func (s *ScaleInstance) Apply(dm *enaml.DeploymentManifest) error {
	ig := dm.GetInstanceGroupByName(s.InstanceGroup)
	///_ = ig // to remove unused error.

	ig.Instances = s.Scale

	return nil
}

func (s *ScaleInstance) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("scale", flag.ContinueOnError)
	fs.StringVar(&s.InstanceGroup, "instance-group", "", "name of the instance group")
	fs.IntVar(&s.Scale, "instances", -1, "number of instances")

	return fs
}

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

	/*	if s.Scale ==  {
			return nil, errors.New("Missing required flag -instance-group")
		}
	*/
	return s, err
}
