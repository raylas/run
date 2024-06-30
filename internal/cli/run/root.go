package run

import (
	"github.com/spf13/cobra"
)

var rootCmdFlags struct {
	attach  bool
	bind    bool
	local   bool
	secrets []string
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run autosemantic scripts",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.attach, "attach", "a", false, "Attach to running script")
	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.bind, "bind", "b", false, "Bind pod to host network")
	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.local, "local", "l", false, "Run script locally")
	cmd.PersistentFlags().StringSliceVarP(&rootCmdFlags.secrets, "secret", "s", []string{}, "Secrets to mount into environment")

	return cmd
}
