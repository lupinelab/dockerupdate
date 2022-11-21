package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// switch args[0] {
		// case "bash":
		dockerupdateBashCompletion(os.Stdout)
		// case "zsh":
		// 	cmd.Root().GenZshCompletion(os.Stdout)
		// case "fish":
		// 	cmd.Root().GenFishCompletion(os.Stdout, true)
		// case "powershell":
		// 	cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		// }
	},
}

func dockerupdateBashCompletion(w io.Writer) error {
	bashCompletion := "#/usr/bin/env bash\n\ncomplete -C 'makelist() { ls $HOME/docker/*/; }; filter() { makelist | grep '^$2' | sort -u; }; filter' dockerupdate"
	buf := new(bytes.Buffer)
	if len(bashCompletion) > 0 {
		buf.WriteString(bashCompletion + "\n")
	}
	_, err := buf.WriteTo(w)
	return err
}
