package pkg

import (
	"fmt"
	"os/exec"
)

func UpdateImage(targetDir string) error {
	var output []byte
	pull := exec.Command("docker-compose", "pull")
	pull.Dir = targetDir
	pullOutput, err := pull.Output()
	if err != nil {
		return err
	}
	output = append(output, pullOutput...)
	recreate := exec.Command("docker-compose", "up", "--force-recreate", "-d")
	recreate.Dir = targetDir
	recreateOutput, err := recreate.Output()
	if err != nil {
		return err
	}
	output = append(output, recreateOutput...)
	fmt.Println(output)
	return nil
}
