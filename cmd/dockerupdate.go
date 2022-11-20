package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lupinelab/dockerupdate/internal"
	"github.com/lupinelab/dockerupdate/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Target struct {
	target string
	name   string
}

func (t *Target) updateContainer() {
	err := pkg.UpdateContainer(t.target)
	if err != nil {
		fmt.Println(err)
	}
}

func (t *Target) updateImage() {
	err := pkg.UpdateImage(t.target)
	if err != nil {
		fmt.Println(err)
	}
}

func (t *Target) buildImage() {
	err := pkg.BuildImage(t.target)
	if err != nil {
		fmt.Println(err)
	}
}

func (t *Target) containerStatus() (string, error) {
	status, err := pkg.ContainerStatus(t.target)
	return status, err
}

func NewTarget(target string) *Target {
	name := filepath.Base(target)
	return &Target{target: target, name: name}
}

func init() {
	dockerupdateCmd.PersistentFlags().BoolP("container", "c", false, "Update container")
	dockerupdateCmd.PersistentFlags().BoolP("image", "i", false, "Update image")
	dockerupdateCmd.PersistentFlags().BoolP("build", "b", false, "Build image")
	dockerupdateCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	dockerupdateCmd.PersistentFlags().Lookup("help").Hidden = true
	cobra.EnableCommandSorting = false
	dockerupdateCmd.CompletionOptions.DisableDefaultCmd = true
}

var dockerupdateCmd = &cobra.Command{
	Use:   "dockerupdate CONTAINER/IMAGE",
	Short: "Perform a task on a container/image",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alltargets := false
		if len(args) == 0 || args[0] == "all" {
			newargs, err := internal.AllTargets()
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			args = newargs
			alltargets = true
		}
		targets := args
		flag := checkFlags(cmd.Flags())
		// the gubbins
		for i := range targets {
			composeTarget, err := internal.ValidateArg(targets[i])
			if err != nil {
				fmt.Print(err.Error())
				continue
			}
			target := NewTarget(composeTarget)
			switch flag {
			case "container":
				fmt.Printf(strings.ToUpper(target.name + "\n" + strings.Repeat("=", len(target.name)) + "\n"))
				target.updateContainer()
				status, err := target.containerStatus()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("Starting %s ... %s\n", target.name, status)
				break
			case "image":
				fmt.Printf(strings.ToUpper(target.name + "\n" + strings.Repeat("=", len(target.name)) + "\n"))
				target.updateImage()
				status, err := target.containerStatus()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("Starting %s ... %s\n", target.name, status)
				break
			case "build":
				fmt.Printf(strings.ToUpper(target.name + "\n" + strings.Repeat("=", len(target.name)) + "\n"))
				target.buildImage()
				break
			}
		}
		if alltargets {
			fmt.Println("\nStatus Summary\n==============")
			for i := range targets {
				status, err := pkg.ContainerStatus(targets[i])
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Printf("%s:"+"%s"+"%s\n", targets[i], strings.Repeat(" ", len(targets[i])), status)
			}
		}
	},
}

func Execute() error {
	return dockerupdateCmd.Execute()
}

func checkFlags(cmdflags *pflag.FlagSet) string {
	var flags []string
	cmdflags.Visit(func(f *pflag.Flag) {
		flags = append(flags, f.Name)
	})
	if len(flags) < 1 {
		fmt.Println("Please provide a flag")
		os.Exit(0)
	}
	if len(flags) > 1 {
		fmt.Println("Please provide a single flag")
		os.Exit(0)
	}
	return flags[0]
}
