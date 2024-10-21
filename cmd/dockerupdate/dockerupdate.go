package dockerupdate

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lupinelab/dockerupdate/internal/targets"
	"github.com/spf13/cobra"
)

func init() {
	dockerupdateCmd.PersistentFlags().BoolP("container", "c", false, "Update container(s)")
	dockerupdateCmd.PersistentFlags().BoolP("image", "i", false, "Update image(s)")
	dockerupdateCmd.PersistentFlags().BoolP("build", "b", false, "Build image(s)")
	dockerupdateCmd.PersistentFlags().BoolP("all", "a", false, "All targets")
	dockerupdateCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	dockerupdateCmd.MarkFlagsMutuallyExclusive("container", "image")
	dockerupdateCmd.MarkFlagsMutuallyExclusive("image", "build")
	dockerupdateCmd.PersistentFlags().Lookup("help").Hidden = true
	dockerupdateCmd.SilenceUsage = true
	cobra.EnableCommandSorting = false
}

var dockerupdateCmd = &cobra.Command{
	Use:   "dockerupdate TARGETDIR [CONTAINER]...",
	Short: "Perform a docker-compose task on a container/image",
	Long: `Perform a docker compose task using a docker-compose file in the TARGETDIR or each TARGETDIR/CONTAINER directory. 
No CONTAINER argument is required if the "all" flag is passed, each subdirectory of the TARGETDIR will be processed.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("all").Changed {
			return cobra.MinimumNArgs(1)(cmd, args)
		}

		return cobra.MinimumNArgs(2)(cmd, args)
	},
	PreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Flag("all").Changed {
			cmd.MarkFlagsOneRequired("build", "container", "image")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		//  Check for docker-compose binary in $PATH
		_, err := exec.LookPath("docker-compose")
		if err != nil {
			return err
		}
		// Assign and validate targets
		targetDir, err := targets.TargetDir(args[0])
		if err != nil {
			return err
		}
		
		var allTargets []string
		allTargets = append(allTargets, targetDir)
		if len(args) > 1 {
			allTargets = args[1:]
		}

		if cmd.Flag("all").Changed {
			// Get all potential targets
			potTargets, err := targets.AllTargets(targetDir)
			if err != nil {
				return err
			}
			if len(allTargets) == 0 {
				return fmt.Errorf("Error: No valid target directories in %s\n", targetDir)
			}
			allTargets = potTargets
		}
		// Make target list
		var targetList []targets.Target
		for _, target := range allTargets {
			composeTarget, err := targets.ValidateArg(targetDir, target)
			if err != nil {
				fmt.Println(err)
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
					fmt.Println(err)
					continue
				}
			}
			if cmd.Flag("image").Changed {
				err := target.UpdateImage()
				if err != nil {
					fmt.Println(err)
					continue
				}
				status, err := target.ContainerStatus()
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("Starting %s ... %s\n", target.Name, status)
			}
			if cmd.Flag("container").Changed {
				err := target.UpdateContainer()
				if err != nil {
					fmt.Println(err)
					continue
				}
				status, err := target.ContainerStatus()
				if err != nil {
					fmt.Println(err)
					continue
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
					fmt.Println(err)
					continue
				}
				fmt.Printf("%s:"+"%s"+"%s\n", target.Name, strings.Repeat(" ", (30-len(target.Name))), status)
			}
		}

		return nil
	},
}

func Execute() error {
	return dockerupdateCmd.Execute()
}
