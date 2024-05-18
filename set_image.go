package composeedit

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/theoremoon/compose-edit/compose"
	cli "github.com/urfave/cli/v2"
)

func setImageCommand() *cli.Command {
	return &cli.Command{
		Name: "set-image",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "image-prefix",
			},
		},
		Action: func(c *cli.Context) error {
			composeFile := c.Args().Get(0)
			if composeFile == "" {
				return errors.New("argument missing: compose file")
			}
			prefix := c.String("image-prefix")
			if prefix == "" {
				return errors.New("flag missing: image-prefix")
			}

			var err error
			prefix, err = normalizeImagePrefix(prefix)
			if err != nil {
				return err
			}

			config, err := compose.LoadFromFile(composeFile)
			if err != nil {
				return err
			}

			// プロパティを上書きしたいのでindexアクセスする
			for i, _ := range config.Services {
				if config.Services[i].Image != "" {
					continue
				}
				config.Services[i].Image = prefix + config.Services[i].Name
			}

			yaml, err := config.MarshalYAML()
			if err != nil {
				return err
			}
			fmt.Println(string(yaml))

			return nil
		},
	}
}

func normalizeImagePrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", errors.New("image prefix is empty")
	}
	if prefix[len(prefix)-1] != '/' {
		return prefix + "/", nil
	}
	return prefix, nil
}
