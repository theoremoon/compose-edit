package composeedit

import (
	"fmt"
	"os"

	composecli "github.com/compose-spec/compose-go/v2/cli"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func setImageCommand() *cli.Command {
	return &cli.Command{
		Name: "set-image",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "image-prefix",
			},
			&cli.BoolFlag{
				Name:  "i",
				Value: false,
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
			prefix := c.String("image-prefix")
			if prefix == "" {
				return errors.New("flag missing: image-prefix")
			}

			prefix, err = normalizeImagePrefix(prefix)
			if err != nil {
				return err
			}

			// プロパティを上書きしたいのでindexアクセスする
			for i, _ := range p.Services {
				if p.Services[i].Image != "" {
					continue
				}
				// avoid Unaddressable Field Assign
				svc := p.Services[i]
				svc.Image = prefix + svc.Name
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

func normalizeImagePrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", errors.New("image prefix is empty")
	}
	if prefix[len(prefix)-1] != '/' {
		return prefix + "/", nil
	}
	return prefix, nil
}
