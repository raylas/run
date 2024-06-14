package cli

import (
	"os"

	"github.com/linecard/run/catalog"
	"github.com/linecard/run/internal/cli/run"
	"github.com/linecard/run/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = run.NewRootCmd()

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	scripts, err := catalog.List()
	if err != nil {
		os.Exit(1)
	}

	for _, script := range scripts {
		rootCmd.AddCommand(run.NewScriptCmd(script))
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Version = config.Version + " (" + config.Commit + ")"
}

func initConfig() {
	config.LoadDefaults()

	viper.AutomaticEnv()
}
