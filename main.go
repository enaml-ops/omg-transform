package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/enaml-ops/enaml"
	"gopkg.in/urfave/cli.v2"
	yaml "gopkg.in/yaml.v2"
)

var Version = "v0.0.1"

func main() {
	app := cli.App{
		Version: Version,
		Flags:   []cli.Flag{
		// TODO: flag for reading input from file
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "change-network",
				Aliases:     []string{"cn"},
				Description: "change the network of an instance group",
				Action:      changeNetwork,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "instance-group",
						Aliases: []string{"i"},
						Usage:   "the instance group whose network you wish to change",
					},
					&cli.StringFlag{
						Name:    "network",
						Aliases: []string{"n"},
						Usage:   "the name of the new network to use",
					},
				},
			},
			&cli.Command{
				Name:        "clone",
				Description: "clone an existing instance group",
				Action:      cloneInstanceGroup,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "instance-group",
						Aliases: []string{"i"},
						Usage:   "the instance group whose network you wish to change",
					},
					&cli.StringFlag{
						Name:    "clone",
						Aliases: []string{"c"},
						Usage:   "the name to use for the clone",
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func loadManifest(c *cli.Context) (*enaml.DeploymentManifest, error) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	manifest := enaml.NewDeploymentManifest(b)
	if manifest == nil {
		return nil, fmt.Errorf("invalid input manifest")
	}
	return manifest, nil
}

func changeNetwork(c *cli.Context) error {
	manifest, err := loadManifest(c)
	if err != nil {
		return err
	}

	groupName := c.String("instance-group")
	if groupName == "" {
		return errors.New("missing required --instance-group flag")
	}

	nc := NetworkMover{
		InstanceGroup: groupName,
		Network:       c.String("network"),
	}
	err = nc.Apply(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	b, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	_, err = os.Stdout.Write(b)
	return err
}

func cloneInstanceGroup(c *cli.Context) error {
	manifest, err := loadManifest(c)
	if err != nil {
		return err
	}

	groupName := c.String("instance-group")
	if groupName == "" {
		return errors.New("missing required --instance-group flag")
	}

	cloner := Cloner{
		InstanceGroup: groupName,
		Clone:         c.String("clone"),
	}
	err = cloner.Apply(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	b, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	_, err = os.Stdout.Write(b)
	return err
}
