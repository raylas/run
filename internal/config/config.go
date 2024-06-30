package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	Version = "dev"
	Commit  = "none"
)

func LoadDefaults() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// run
	viper.SetDefault("image", "ghcr.io/nixos/nix")
	viper.SetDefault("entrypoint", []string{"bash", "-c"})
	viper.SetDefault("command", "echo %s | base64 -d -i > /run && chmod +x /run && /run %s")

	// Kubernetes
	viper.SetDefault("kubernetes.config_path", homeDir+"/.kube/config")
	viper.SetDefault("kubernetes.namespace", "default")
}
