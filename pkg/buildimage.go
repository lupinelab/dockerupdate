package pkg

import (
	"fmt"
	"os/exec"
)

func BuildImage(targetDir string) error {
	build := exec.Command("docker-compose", "build")
	build.Dir = targetDir
	output, err := build.Output()
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
