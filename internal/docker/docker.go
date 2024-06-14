package docker

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

const (
	image = "ghcr.io/nixos/nix"
)

func Run(ctx context.Context, attach bool, cmd string) (string, error) {
	c, err := client.NewClientWithOpts(
		client.WithHost(viper.GetString("docker.host")),
		client.WithVersion(viper.GetString("docker.api.version")),
	)
	if err != nil {
		return "", err
	}

	ctr, err := c.ContainerCreate(
		ctx,
		&container.Config{
			Image:      image,
			Entrypoint: []string{"bash", "-c"},
			Cmd:        []string{cmd},
			Tty:        attach,
			OpenStdin:  attach,
			StdinOnce:  attach,
		},
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", err
	}

	if err := c.ContainerStart(ctx, ctr.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	switch {
	case attach:
		if err := Attach(ctx, c, ctr.ID); err != nil {
			return "", err
		}
	default:
		return ctr.ID, nil
	}

	return "", nil
}

func Attach(ctx context.Context, c *client.Client, id string) error {
	attachment, err := c.ContainerAttach(ctx, id, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer attachment.Close()

	prevTerm, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), prevTerm)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGWINCH)

	resizeTty := func() {
		if w, h, err := term.GetSize(0); err == nil {
			c.ContainerResize(ctx, id, container.ResizeOptions{
				Height: uint(h),
				Width:  uint(w),
			})
		}
	}

	resizeTty()
	go func() {
		for range sigCh {
			resizeTty()
		}
	}()

	// stdout
	go func() {
		io.Copy(os.Stdout, attachment.Reader)
		sigCh <- syscall.SIGTERM
	}()

	// stdin
	go func() {
		io.Copy(attachment.Conn, os.Stdin)
		sigCh <- syscall.SIGTERM
	}()

	<-sigCh

	if err := c.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		return err
	}

	return nil
}
