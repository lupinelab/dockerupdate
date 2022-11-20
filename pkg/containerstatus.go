package pkg

import (
	"context"
	"path/filepath"

	"github.com/docker/docker/client"
)

func ContainerStatus(targetDir string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// Replace this ID with a container that really exists
	out, err := cli.ContainerInspect(ctx, filepath.Base(targetDir))
	if err != nil {
		return "", err
	}

	return out.State.Status, nil
}
