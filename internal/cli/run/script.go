package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/linecard/run/catalog"
	"github.com/linecard/run/internal/docker"
	"github.com/linecard/run/internal/equip"
	"github.com/linecard/run/internal/script"
	"github.com/spf13/cobra"
)

func NewScriptCmd(name string) *cobra.Command {
	scriptDesc, scriptArgs, err := script.ParseSpec(name)
	if err != nil {
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   name,
		Short: *scriptDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			var flagArgs = []string{}
			for _, arg := range scriptArgs {
				flagArgs = append(flagArgs, cmd.Flag(arg.Name).Value.String())
			}

			switch {
			case rootCmdFlags.local:
				packed, err := equip.Pack(name, strings.Join(flagArgs, " "))
				if err != nil {
					return err
				}

				ctr, err := docker.Run(cmd.Context(), rootCmdFlags.attach, packed)
				if err != nil {
					return err
				}

				if ctr != "" {
					fmt.Println(ctr)
				}

				return nil
			default:
				return nil
			}
		},
	}

	for _, arg := range scriptArgs {
		cmd.Flags().String(arg.Name, arg.Value, arg.Desc)

		if arg.Value == "" {
			cmd.MarkFlagRequired(arg.Name)
		}
	}

	cmd.AddCommand(NewInspectCmd(name))

	return cmd
}

func NewInspectCmd(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect script",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.NoArgs(cmd, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			script, err := catalog.Read(name)
			if err != nil {
				return err
			}

			fmt.Println(string(script))

			return nil
		},
	}

	return cmd
}
