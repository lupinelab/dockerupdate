package targets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func TargetDir(targetdir string) (string, error) {
	// Check targetdir exists
	targetDirPath, err := filepath.Abs(targetdir)
	if err != nil {
		return targetDirPath, err
	}
	fileInfo, err := os.Stat(targetDirPath)
	if err != nil {
		return targetDirPath, err
	}
	// Check if targetdir is a directory
	if fileInfo.IsDir() == false {
		err := errors.New(fmt.Sprintf("Error: %s is not a directory", targetDirPath))
		return targetDirPath, err
	}
	// Docker dir path
	return targetDirPath, nil
}
