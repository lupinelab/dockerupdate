package dockerupdate

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lupinelab/dockerupdate/internal/targets"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	dockerupdateCmd.PersistentFlags().BoolP("container", "c", false, "Update container")
	dockerupdateCmd.PersistentFlags().BoolP("image", "i", false, "Update image")
	dockerupdateCmd.PersistentFlags().BoolP("build", "b", false, "Build image")
	dockerupdateCmd.PersistentFlags().BoolP("all", "a", false, "All targets")
	dockerupdateCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	dockerupdateCmd.MarkFlagsMutuallyExclusive("container", "image")
	dockerupdateCmd.MarkFlagsMutuallyExclusive("image", "build")
	dockerupdateCmd.PersistentFlags().Lookup("help").Hidden = true
	cobra.EnableCommandSorting = false
}

var dockerupdateCmd = &cobra.Command{
	Use:   "dockerupdate TARGETDIR [CONTAINER]...",
	Short: "Perform a docker-compose task on a container/image",
	Long: `Perform a docker compose task using a docker-compose file in the TARGETDIR or each TARGETDIR/CONTAINER directory. 
No CONTAINER argument is required if the "all" flag is passed, each subdirectory of the TARGETDIR will be processed.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arg and "all" flag has not been set show help
		if len(args) == 0 && !cmd.Flag("all").Changed {
			cmd.Help()
			return
		}
		// Get the flags
		var flags []string
		cmd.Flags().Visit(func(f *pflag.Flag) {
			flags = append(flags, f.Name)
		})
		if len(flags) == 0 {
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
		targetDir, err := targets.TargetDir(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		var allTargets []string
		allTargets = append(allTargets, targetDir)
		if len(args) > 1 {
			allTargets = args[1:]
		}

		if cmd.Flag("all").Changed {
			// Check for another flag if "all" flag is set
			if len(flags) == 1 {
				fmt.Println("Error: if the flag [all] is set one of the following flags must also be set: [build image container]")
				return
			}
			// Get all potential targets
			potTargets, err := targets.AllTargets(targetDir)
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			if len(allTargets) == 0 {
				fmt.Printf("Error: No valid target directories in %s\n", targetDir)
				return
			}
			allTargets = potTargets
		}
		// Make target list
		var targetList []targets.Target
		for _, target := range allTargets {
			composeTarget, err := targets.ValidateArg(targetDir, target)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			target := targets.NewTarget(composeTarget)
			targetList = append(targetList, *target)
		}
		// The gubbins
		for _, target := range targetList {
			fmt.Print(strings.ToUpper(target.Name + "\n" + strings.Repeat("=", len(target.Name)) + "\n"))
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
