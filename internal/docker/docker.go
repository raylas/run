package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/linecard/run/internal/output"
	"github.com/spf13/viper"
)

func Run(ctx context.Context, attach bool, cmd, name string) (string, error) {
	c, err := client.NewClientWithOpts(
		client.WithHostFromEnv(),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return "", err
	}

	ctr, err := c.ContainerCreate(
		ctx,
		&container.Config{
			Image:      viper.GetString("image"),
			Entrypoint: viper.GetStringSlice("entrypoint"),
			Cmd:        []string{cmd},
			Tty:        attach,
			OpenStdin:  attach,
			StdinOnce:  attach,
		},
		nil,
		nil,
		nil,
		name,
	)
	if err != nil {
		return "", err
	}

	if err := c.ContainerStart(ctx, ctr.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	switch {
	case attach:
		if err := tty(ctx, c, ctr.ID); err != nil {
			return "", err
		}
	default:
		return output.Colorize(output.Green, "Container started: %s", ctr.ID[:12]), nil
	}

	return output.Colorize(output.Green, "Container TTY closed"), nil
}
