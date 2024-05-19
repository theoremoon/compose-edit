package composeedit

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func setImageCommand(loadCompose loadComposeFunc) *cli.Command {
	return &cli.Command{
		Name: "set-image",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "format",
			},
			&cli.BoolFlag{
				Name:  "i",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			p, err := loadCompose(c)
			if err != nil {
				return err
			}
			format := c.String("format")
			if format == "" {
				return errors.New("flag missing: image-format")
			}

			// プロパティを上書きしたいのでindexアクセスする
			for i, _ := range p.Services {
				if p.Services[i].Image != "" {
					continue
				}
				// avoid Unaddressable Field Assign
				svc := p.Services[i]
				svc.Image = strings.Replace(format, "{Name}", svc.Name, -1)
				p.Services[i] = svc
			}

			yaml, err := p.MarshalYAML()
			if err != nil {
				return err
			}
			if c.Bool("i") {
				f := p.ComposeFiles[0]
				err = os.WriteFile(f, yaml, 0644)
				if err != nil {
					return err
				}
			} else {
				fmt.Println(string(yaml))
			}

			return nil
		},
	}
}
