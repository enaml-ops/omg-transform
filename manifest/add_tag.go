package manifest

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/enaml-ops/enaml"
)

type TagAdder struct {
	Args []string
}

func (t *TagAdder) Apply(dm *enaml.DeploymentManifest) error {
	for _, arg := range t.Args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		dm.AddTag(key, value)
	}
	return nil
}

func (t *TagAdder) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("add-tag", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of add-tags:")
		fmt.Fprintln(os.Stderr, "   add-tags key=value key=value ...")
	}
	return fs
}

func AddTagsTransformation(args []string) (Transformation, error) {
	t := &TagAdder{}

	// we don't actually need a FlagSet here, but we use it to get a nice
	// usage message if given a `-help` argument.
	fs := t.flagSet()
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, errors.New("missing tag specifier(s) [format key=value]")
	}

	t.Args = fs.Args()
	for _, arg := range t.Args {
		if c := strings.Count(arg, "="); c != 1 {
			return nil, fmt.Errorf("invalid tag specifier %q, expected format key=value", arg)
		}
		parts := strings.Split(arg, "=")
		if parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid tag specifier %q, expected format key=value", arg)
		}
	}
	return t, nil
}
