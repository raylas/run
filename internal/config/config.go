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

	// Docker
	// viper.SetDefault("docker.host", "unix://"+homeDir+".colima/default/docker.sock")
	viper.SetDefault("docker.api.version", "1.43")

	// Kubernetes
	viper.SetDefault("kubernetes.config_path", homeDir+"/.kube/config")
}
