package internal

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/client"
)

type Target struct {
	target string
	Name   string
}

func (t *Target) UpdateContainer() error {
	recreate := exec.Command("docker-compose", "up", "--force-recreate", "--no-build", "-d")
	recreate.Dir = t.target
	recreateOutput, err := recreate.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(recreateOutput[:]))
	return err
}

func (t *Target) UpdateImage() error {
	pull := exec.Command("docker-compose", "build")
	pull.Dir = t.target
	pullOutput, err := pull.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(pullOutput))
	recreate := exec.Command("docker-compose", "up", "--force-recreate", "--no-build", "-d")
	recreate.Dir = t.target
	recreateOutput, err := recreate.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(recreateOutput[:]))
	return err
}

func (t *Target) BuildImage() error {
	build := exec.Command("docker-compose", "build")
	build.Dir = t.target
	buildOutput, err := build.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(buildOutput[:]))
	return err
}

func (t *Target) ContainerStatus() (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()
	out, err := cli.ContainerInspect(ctx, filepath.Base(t.target))
	if err != nil {
		return "", err
	}

	return out.State.Status, nil
}

func NewTarget(target string) *Target {
	name := filepath.Base(target)
	return &Target{target: target, Name: name}
}
