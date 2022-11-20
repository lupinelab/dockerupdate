package cmd

import (
	"fmt"
	"os"

	"github.com/lupinelab/dockerupdate/internal"
	"github.com/lupinelab/dockerupdate/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Target struct {
	target string
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
		// Check arg(s)
		if len(args) < 1 {
			fmt.Println("Please provide a container/image name")
			return
		}
		flag := checkFlags(cmd.Flags())
		// the gubbins
		for i := range args {
			composeTarget, err := internal.ValidateArg(args[i])
			if err != nil {
				fmt.Print(err.Error())
			}
			target := Target{composeTarget}
			switch flag {
			case "container":
				target.updateContainer()
				break
			case "image":
				target.updateImage()
				break
			case "build":
				target.buildImage()
				break
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
