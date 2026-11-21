package cli

import (
	"fmt"
	"log"
	"os"

	"go.rmdd.net/run/catalog"
	runcmd "go.rmdd.net/run/internal/cli/run"
	"go.rmdd.net/run/internal/config"
	"go.rmdd.net/run/internal/output"
	"go.rmdd.net/run/internal/parse"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = runcmd.NewRootCmd()

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.LoadDefaults)

	viper.SetEnvPrefix("run")
	viper.AutomaticEnv()

	scripts, err := catalog.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range scripts {
		contents, err := catalog.Read(s)
		if err != nil {
			fmt.Println(output.Colorize(output.Red, "%v", err))
			continue
		}

		desc, args, err := parse.Args(contents)
		if err != nil {
			fmt.Println(output.Colorize(output.Red, "%v", err))
			continue
		}

		rootCmd.AddCommand(runcmd.NewScriptCmd(s, desc, args))
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Version = config.Version + " (" + config.Commit + ")"
}
