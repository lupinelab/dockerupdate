package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lupinelab/dockerupdate/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	dockerupdateCmd.PersistentFlags().BoolP("container", "c", false, "update container")
	dockerupdateCmd.PersistentFlags().BoolP("image", "i", false, "update image")
	dockerupdateCmd.PersistentFlags().BoolP("build", "b", false, "build image")
	dockerupdateCmd.PersistentFlags().BoolP("all", "a", false, "all targets")
	dockerupdateCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	dockerupdateCmd.MarkFlagsMutuallyExclusive("container", "image")
	dockerupdateCmd.MarkFlagsMutuallyExclusive("image", "build")
	dockerupdateCmd.PersistentFlags().Lookup("help").Hidden = true
	dockerupdateCmd.AddCommand(completionCmd)
	dockerupdateCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.EnableCommandSorting = false
}

var dockerupdateCmd = &cobra.Command{
	Use:   "dockerupdate CONTAINER/IMAGE",
	Short: "Perform a docker-compose task on a container/image",
	Long: `Perform a docker compose task on a container/image in the $HOME/docker/$1 directory. 
No arguement is required if the "all" flag is passed, all directories in $HOME/docker will be processed.
If the docker-compose binary is not in $PATH an error will be returned.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arg and "all" flag has not been setsudo show help
		if len(args) == 0 && !cmd.Flag("all").Changed {
			cmd.Help()
			return
		}
		//  Check for docker-compose binary in $PATH
		_, err := exec.LookPath("docker-compose")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// Assign and validate targets
		targets := args
		if cmd.Flag("all").Changed {
			// Check for another flag if "all" flag is set
			var flags []string
			cmd.Flags().Visit(func(f *pflag.Flag) {
				flags = append(flags, f.Name)
			})
			if len(flags) == 1 {
				fmt.Println("Error: if the flag [all] is set one of the following flags must also be set: [build image container]")
				return
			}
			// Get all potential targets
			allTargets, err := internal.AllTargets()
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			targets = allTargets
		}
		// Make target list
		var targetList []internal.Target
		for i := range targets {
			composeTarget, err := internal.ValidateArg(targets[i])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			target := internal.NewTarget(composeTarget)
			targetList = append(targetList, *target)
		}
		// The gubbins
		for _, target := range targetList {
			fmt.Printf(strings.ToUpper(target.Name + "\n" + strings.Repeat("=", len(target.Name)) + "\n"))
			if cmd.Flag("build").Changed {
				err := target.BuildImage()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			if cmd.Flag("image").Changed {
				err := target.UpdateImage()
				if err != nil {
					fmt.Println(err.Error())
				}
				status, err := target.ContainerStatus()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("Starting %s ... %s\n", target.Name, status)
			}
			if cmd.Flag("container").Changed {
				err := target.UpdateContainer()
				if err != nil {
					fmt.Println(err.Error())
				}
				status, err := target.ContainerStatus()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("Starting %s ... %s\n", target.Name, status)
			}
		}
		// Print summary if "all" flag set
		if cmd.Flag("all").Changed {
			fmt.Println("\nStatus Summary\n==============")
			for _, target := range targetList {
				status, err := target.ContainerStatus()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Printf("%s:"+"%s"+"%s\n", target.Name, strings.Repeat(" ", (30-len(target.Name))), status)
			}
		}
	},
}

func Execute() error {
	return dockerupdateCmd.Execute()
}
