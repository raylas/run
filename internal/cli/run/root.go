package run

import (
	"github.com/raylas/run/catalog"
	"github.com/spf13/cobra"
)

var rootCmdFlags struct {
	attach     bool
	bind       bool
	capture    bool
	local      bool
	secretEnv  []string
	secretFile []string
	clearCache bool
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run autosemantic scripts",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				// permit -x/--clear-cache by itself, otherwise require exactly one arg
				if rootCmdFlags.clearCache && len(args) == 0 {
					return nil
				}
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we need to clear the cache
			if rootCmdFlags.clearCache {
				if err := catalog.ClearCache(); err != nil {
					return err
				}
				return nil
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.attach, "attach", "a", false, "Attach to running script")
	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.bind, "bind", "b", false, "Bind pod to host network")
	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.capture, "capture", "c", false, "Enable packet capture (implies --bind)")
	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.local, "local", "l", false, "Run script locally")
	cmd.PersistentFlags().StringSliceVarP(&rootCmdFlags.secretEnv, "secret-env", "s", []string{}, "Secrets to mount into environment")
	cmd.PersistentFlags().StringSliceVarP(&rootCmdFlags.secretFile, "secret-file", "f", []string{}, "Secrets to mount into file")
	cmd.PersistentFlags().BoolVarP(&rootCmdFlags.clearCache, "clear-cache", "x", false, "Clear the remote catalog cache")

	return cmd
}
