package composeedit

import (
	"fmt"
	"regexp"

	composecli "github.com/compose-spec/compose-go/v2/cli"
	cli "github.com/urfave/cli/v2"
)

func listImagesCommand() *cli.Command {
	return &cli.Command{
		Name: "list-images",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "filter",
			},
		},
		Action: func(c *cli.Context) error {
			o, err := composecli.NewProjectOptions([]string{}, composecli.WithDefaultConfigPath)
			if err != nil {
				return err
			}
			p, err := o.LoadProject(c.Context)
			if err != nil {
				return err
			}
			filter := c.String("filter")
			if filter == "" {
				filter = ".*"
			}
			pat, err := regexp.Compile(filter)
			if err != nil {
				return err
			}

			for _, svc := range p.Services {
				if svc.Image == "" {
					continue
				}
				if pat.MatchString(svc.Image) {
					fmt.Println(svc.Image)
				}

			}

			return nil
		},
	}
}
