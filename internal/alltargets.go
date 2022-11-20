package internal

import (
	"os"
)

func AllTargets() ([]string, error) {
	var allTargets []string
	f, err := os.Open(DockerDir())
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}
	for _, p := range files {
		if p.IsDir() {
			allTargets = append(allTargets, p.Name())
		}
	}
	return allTargets, nil
}
