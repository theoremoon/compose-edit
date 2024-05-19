package composeedit

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	composelib "github.com/compose-spec/compose-go/v2/types"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func taskifyCommand(loadCompose loadComposeFunc) *cli.Command {
	return &cli.Command{
		Name: "taskify",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "family",
			},
			&cli.StringFlag{
				Name: "execution-role",
			},
		},
		Action: func(c *cli.Context) error {
			p, err := loadCompose(c)
			if err != nil {
				return err
			}
			family := c.String("family")
			if family == "" {
				return errors.New("family is required")
			}
			executionRole := c.String("execution-role")
			if executionRole == "" {
				return errors.New("execution-role is required")
			}

			td, _, err := taskify(p, family, executionRole)
			if err != nil {
				return err
			}

			tdjson, err := MarshalJSONForAPI(td, "del(.ipcMode)")
			if err != nil {
				return err
			}
			fmt.Println(string(tdjson))

			return nil
		},
	}
}

func taskify(p *composelib.Project, family, executionRole string) (*ecs.RegisterTaskDefinitionInput, *ecs.CreateServiceInput, error) {
	cdefs := make([]types.ContainerDefinition, 0)
	for _, svc := range p.Services {
		cdef, err := toContainerDefinition(&svc)
		if err != nil {
			return nil, nil, err
		}
		cdefs = append(cdefs, *cdef)
	}
	td := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: cdefs,
		Family:               &family,
		// XXX: dummy values for now. it should be configurable
		Cpu:                     aws.String("256"),
		Memory:                  aws.String("512"),
		NetworkMode:             types.NetworkModeAwsvpc,
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		PidMode:                 types.PidModeTask, // for Fargate, it should be task
		ExecutionRoleArn:        &executionRole,
	}

	return td, nil, nil
}

func toContainerDefinition(svc *composelib.ServiceConfig) (*types.ContainerDefinition, error) {
	envs := make([]types.KeyValuePair, 0)
	for k, v := range svc.Environment {
		envs = append(envs, types.KeyValuePair{
			Name:  &k,
			Value: v,
		})
	}

	ports := make([]types.PortMapping, 0)
	for _, p := range svc.Ports {

		port := int32(p.Target)
		p64, err := strconv.ParseInt(p.Published, 10, 32)
		if err != nil {
			return nil, err
		}
		published := int32(p64)

		ports = append(ports, types.PortMapping{
			ContainerPort: &port,
			HostPort:      &published,
			AppProtocol:   types.ApplicationProtocolHttp, // よくわからん
			Protocol:      types.TransportProtocol(p.Protocol),
		})
	}

	depends := make([]types.ContainerDependency, 0)
	for name, d := range svc.DependsOn {
		condition, err := convertCondition(d.Condition)
		if err != nil {
			return nil, err
		}
		depends = append(depends, types.ContainerDependency{
			Condition:     condition,
			ContainerName: &name,
		})
	}

	cdef := types.ContainerDefinition{
		Name:         &svc.Name,
		Image:        &svc.Image,
		Command:      svc.Command,
		EntryPoint:   svc.Entrypoint,
		Environment:  envs,
		PortMappings: ports,
		Essential:    aws.Bool(true),
		DependsOn:    depends,
	}
	return &cdef, nil
}

func convertCondition(c string) (types.ContainerCondition, error) {
	switch c {
	case "service_healthy":
		return types.ContainerCondition("HEALTHY"), nil
	case "service_started":
		return types.ContainerCondition("START"), nil
	case "service_completed_successfully":
		return types.ContainerCondition("SUCCESS"), nil
	default:
		return "", errors.New("unknown condition: " + c)
	}
}
