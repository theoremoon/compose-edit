package composeedit

import (
	cli "github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name: "compose-edit",
		Commands: []*cli.Command{
			verifyCommand(),
			setImageCommand(),
			listImagesCommand(),
		},
	}
}
