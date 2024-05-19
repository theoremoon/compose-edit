package composeedit

import (
	"fmt"
	"regexp"

	cli "github.com/urfave/cli/v2"
)

func listImagesCommand(loadCompose loadComposeFunc) *cli.Command {
	return &cli.Command{
		Name: "list-images",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "filter",
			},
		},
		Action: func(c *cli.Context) error {
			p, err := loadCompose(c)
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
