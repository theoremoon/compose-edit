package composeedit

import (
	composecli "github.com/compose-spec/compose-go/v2/cli"
	composelib "github.com/compose-spec/compose-go/v2/types"
	cli "github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name: "compose-edit",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name: "f",
			},
		},
		Commands: []*cli.Command{
			catCommand(loadComposeProject),
			listImagesCommand(loadComposeProject),
			setImageCommand(loadComposeProject),
			verifyCommand(loadComposeProject),
		},
	}
}

type loadComposeFunc func(c *cli.Context) (*composelib.Project, error)

func loadComposeProject(c *cli.Context) (*composelib.Project, error) {
	configPaths := c.StringSlice("f")
	o, err := composecli.NewProjectOptions(configPaths, composecli.WithDefaultConfigPath)
	if err != nil {
		return nil, err
	}
	p, err := o.LoadProject(c.Context)
	if err != nil {
		return nil, err
	}

	return p, nil
}
