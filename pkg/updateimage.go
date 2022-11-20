package pkg

import (
	"fmt"
	"os/exec"
)

func UpdateImage(targetDir string) error {
	var output []byte
	pullcmd := "docker-compose pull"
	pull := exec.Command("bash", "-c", pullcmd)
	pull.Dir = targetDir
	pullOutput, err := pull.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(pullOutput))
	recreate := exec.Command("docker-compose", "up", "--force-recreate", "-d")
	recreate.Dir = targetDir
	recreateOutput, err := recreate.CombinedOutput()
	if err != nil {
		return err
	}
	output = append(output, recreateOutput...)
	fmt.Printf(string(recreateOutput[:]))
	return nil
}
