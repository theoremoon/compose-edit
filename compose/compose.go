package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/loader"
	composelib "github.com/compose-spec/compose-go/types"
)

func LoadFromFile(path string) (*composelib.Project, error) {
	composeBuf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	conf, err := loadCompose(composeBuf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func loadCompose(buf []byte) (*composelib.Project, error) {
	composeConfig, err := loader.ParseYAML(buf)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	dir := filepath.Dir(".")
	config, err := loader.Load(composelib.ConfigDetails{
		WorkingDir: dir,
		ConfigFiles: []composelib.ConfigFile{
			{
				Filename: "compose.yml",
				Config:   composeConfig,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	return config, nil
}
