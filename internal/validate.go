package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func ValidateArg(arg string) (string, error) {
	targetDir := filepath.Join(DockerDir(), arg)
	// Check arg is in dockerDir folder
	fileInfo, err := os.Stat(targetDir)
	if err != nil {
		return "", err
	}
	// Check if arg is a directory
	if fileInfo.IsDir() == false {
		err := fmt.Errorf("Error: %s is not a directory", targetDir)
		return "", err
	}
	// if targetDir.IsDir() {
	// 	fmt.Printf("Error: %s is not a directory", targetDir)
	// }
	// Check for docker-compose.yml or docker-compose.yaml files in targetDir
	f, err := os.Open(targetDir)
	if err != nil {
		return "", err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return "", err
	}
	for _, p := range files {
		if p.IsDir() == false {
			if p.Name() == "docker-compose.yml" || p.Name() == "docker-compose.yaml" {
				composeTarget := targetDir
				return composeTarget, err
			}
		}
	}
	err = fmt.Errorf("Error: No valid docker-compose files found in %s", targetDir)
	return "", err
}
