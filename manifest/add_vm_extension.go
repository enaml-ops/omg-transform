package manifest

import (
	"errors"
	"flag"
	"fmt"

	"github.com/enaml-ops/enaml"
)

// VMExtension is a transformation that adds a vm extension to the given instance group
type VMExtension struct {
	Name          string
	InstanceGroup string
	Extensions    []string
}

func (ve *VMExtension) Apply(dm *enaml.DeploymentManifest) error {
	ig := dm.GetInstanceGroupByName(ve.InstanceGroup)
	if ig == nil {
		return fmt.Errorf("couldn't find instance group %s", ve.InstanceGroup)
	}
	ig.VMExtensions = append(ig.VMExtensions, ve.Extensions...)
	return nil
}

func (ve *VMExtension) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("add-vm-extension", flag.ContinueOnError)
	fs.StringVar(&ve.InstanceGroup, "instance-group", "", "Name of the instance group")
	fs.StringVar(&ve.Name, "name", "", "Name(s) of the vm extension [If multiple, comma separate values]")
	return fs
}

// AddVMExtensionTransformation is a TransformationBuilder that builds the
// 'add-vm-extension' transformation.
func AddVMExtensionTransformation(args []string) (Transformation, error) {
	ve := &VMExtension{}
	fs := ve.flagSet()
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if ve.InstanceGroup == "" {
		return nil, errors.New("missing required flag instance-group")
	}
	if ve.Name == "" {
		return nil, errors.New("missing required flag name")
	}
	ve.Extensions = split(ve.Name, ",")
	if len(ve.Extensions) == 0 {
		return nil, errors.New("invalid format for extension names, must be comma-separated")
	}

	return ve, nil
}
