package internal

import (
	"os"
)

func AllTargets(targetdir string) ([]string, error) {
	var allTargets []string
	targetDir, err := TargetDir(targetdir)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(targetDir)
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
