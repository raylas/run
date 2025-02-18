package catalog

import (
	"embed"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed *
var embeddedCatalog embed.FS

func List() ([]string, error) {
	// Get embedded scripts
	files, err := embeddedCatalog.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var scripts []string
	for _, file := range files {
		if file.Name() == "catalog.go" || file.Name() == "remote.go" || file.IsDir() {
			continue
		}
		scripts = append(scripts, file.Name())
	}

	// If remote catalog URL is set, fetch and add those scripts
	//
	// RUN_REMOTE_CATALOG_URL should be a Go-style package reference, e.g.:
	// github.com/username/repo
	if viper.GetString("remote_catalog_url") != "" {
		fmt.Printf("Fetching remote catalog from: %s", viper.GetString("remote_catalog_url"))
		remoteScripts, err := ListRemote()
		if err != nil {
			return scripts, fmt.Errorf("error fetching remote catalog: %w", err)
		}
		scripts = append(scripts, remoteScripts...)
	}

	return scripts, nil
}

func Read(name string) ([]byte, error) {
	// Try embedded first
	bytes, err := embeddedCatalog.ReadFile(name)
	if err == nil {
		return bytes, nil
	}

	// If not found in embedded and we have a remote URL, try remote
	if viper.GetString("remote_catalog_url") != "" {
		bytes, err := ReadRemote(name)
		if err != nil {
			return nil, fmt.Errorf("script not found in embedded or remote catalog: %s", name)
		}
		return bytes, nil
	}

	return nil, fmt.Errorf("script not found: %s", name)
}
