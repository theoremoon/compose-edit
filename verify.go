package composeedit

import (
	composecli "github.com/compose-spec/compose-go/v2/cli"
	composelib "github.com/compose-spec/compose-go/v2/types"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func verifyCommand() *cli.Command {
	return &cli.Command{
		Name: "verify",
		Action: func(c *cli.Context) error {
			o, err := composecli.NewProjectOptions([]string{}, composecli.WithDefaultConfigPath)
			if err != nil {
				return err
			}
			p, err := o.LoadProject(c.Context)
			if err != nil {
				return err
			}
			if err := ValidateComposeConfig(p); err != nil {
				return err
			}

			return nil
		},
	}
}

func ValidateComposeConfig(conf *composelib.Project) error {
	if len(conf.Volumes) != 0 {
		return errors.New("volume is prohibited")
	}
	for _, svc := range conf.Services {
		if len(svc.CapAdd) != 0 {
			return errors.New("cap_add is prohibited")
		}
		if len(svc.Configs) != 0 {
			return errors.New("configs is prohibited")
		}
		if svc.CredentialSpec != nil {
			return errors.New("credential_spec is prohibited")
		}
		if svc.Deploy != nil {
			return errors.New("deploy is prohibited")
		}
		if len(svc.EnvFiles) != 0 {
			return errors.New("envfile is prohibited")
		}
		if svc.Logging != nil {
			return errors.New("logging is prohibited")
		}
		if svc.OomKillDisable {
			return errors.New("oom_kill_disable is prohibited")
		}
		if svc.Privileged {
			return errors.New("privileged is prohibited")
		}
		if svc.Runtime != "" {
			return errors.New("runtime is prohibited")
		}
		if len(svc.Secrets) != 0 {
			return errors.New("secrets is prohibited")
		}
		if len(svc.Sysctls) != 0 {
			return errors.New("sysctls is prohibited")
		}
		if len(svc.Tmpfs) != 0 {
			return errors.New("tmpfs is prohibited")
		}
		if len(svc.Volumes) != 0 {
			return errors.New("volumes is prohibited")
		}
		if len(svc.VolumesFrom) != 0 {
			return errors.New("volumes_from is prohibited")
		}
		if len(svc.Extensions) != 0 {
			return errors.New("extensions is prohibited")
		}
	}
	return nil
}
