package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func DockerDir() string {
	// User's home dir
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	// Check $HOME/docker exists
	dockerDirPath, err := filepath.Abs(filepath.Join(userHome, "/docker"))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(0)
	}
	fileInfo, err := os.Stat(dockerDirPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	// Check if $HOME/docker is a directory
	if fileInfo.IsDir() == false {
		fmt.Println("Error: $HOME/docker is not a directory")
		os.Exit(0)
	}
	// Docker dir path
	return (dockerDirPath)
}
