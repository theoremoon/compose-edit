package composeedit

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

func catCommand(loadCompose loadComposeFunc) *cli.Command {
	return &cli.Command{
		Name: "cat",
		Action: func(c *cli.Context) error {
			p, err := loadCompose(c)
			if err != nil {
				return err
			}
			yaml, err := p.MarshalYAML()
			if err != nil {
				return err
			}

			fmt.Println(string(yaml))

			return nil
		},
	}
}
