package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ValidateArg(targetdir string, target string) (string, error) {
	if targetdir == target {
		target = filepath.Base(targetdir)
		targetdir = strings.Split(targetdir, target)[0]
	}
	targetAbs := filepath.Join(targetdir, target)
	// Check arg is in dockerDir folder
	fileInfo, err := os.Stat(targetAbs)
	if err != nil {
		return "", err
	}
	// Check if arg is a directory
	if fileInfo.IsDir() == false {
		err := fmt.Errorf("Error: %s is not a directory", targetAbs)
		return "", err
	}
	// Check for docker-compose.yml or docker-compose.yaml files in targetDir and return absolute path to
	f, err := os.Open(targetAbs)
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
				composeTarget := targetAbs
				return composeTarget, err
			}
		}
	}
	err = fmt.Errorf("Error: No valid docker-compose files found in %s\n", targetAbs)
	return "", err
}
