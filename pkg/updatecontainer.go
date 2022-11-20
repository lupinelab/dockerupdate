package pkg

import (
	"fmt"
	"os/exec"
)

func UpdateContainer(targetDir string) error {
	fmt.Println("test")
	recreate := exec.Command("docker-compose", "up", "--force-recreate", "-d")
	recreate.Dir = targetDir
	recreateOutput, err := recreate.Output()
	if err != nil {
		return err
	}
	fmt.Printf(string(recreateOutput[:]))
	return nil
}
