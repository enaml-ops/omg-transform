package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
)

type AZChanger struct {
	InstanceGroup string
	AZs           []string
	azsFlag       string
}

func (a *AZChanger) Apply(dm *enaml.DeploymentManifest) error {
	ig := dm.GetInstanceGroupByName(a.InstanceGroup)
	if ig == nil {
		return fmt.Errorf("couldn't find instance group %s", a.InstanceGroup)
	}

	ig.AZs = a.AZs
	return nil
}

func (a *AZChanger) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("change-az", flag.ExitOnError)
	fs.StringVar(&a.InstanceGroup, "instance-group", "", "name of the instance group")
	fs.StringVar(&a.azsFlag, "az", "", "a comma separated list of az names")
	return fs
}

func ChangeAZTransformation(args []string) (Transformation, error) {
	a := &AZChanger{}
	fs := a.flagSet()
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if a.InstanceGroup == "" {
		return nil, errors.New("missing required flag -instance-group")
	}
	if a.azsFlag == "" {
		return nil, errors.New("missing required flag -az")
	}
	tempazs := strings.Split(a.azsFlag, ",")

	if strings.Contains(a.azsFlag, " ") {
		return nil, errors.New("invalid format for az, cannot contain space")
	}

	for i := range tempazs {
		if tempazs[i] != "" {
			a.AZs = append(a.AZs, tempazs[i])
		}
	}

	if len(a.AZs) == 0 {
		return nil, errors.New("invalid format for az, must be comma-separated")
	}

	return a, nil
}
