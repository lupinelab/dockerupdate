package pkg

import (
	"fmt"
	"os/exec"
)

func BuildImage(targetDir string) error {
	build := exec.Command("docker-compose", "build")
	build.Dir = targetDir
	buildOutput, err := build.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(buildOutput[:]))
	return nil
}
