package catalog

import (
	"embed"
	"fmt"
	"log"
)

//go:embed *
var embeddedCatalog embed.FS

// RemoteCatalogURL should be a Go-style package reference, e.g.:
// github.com/username/repo
var RemoteCatalogURL string

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
	if RemoteCatalogURL != "" {
		log.Printf("Fetching remote catalog from: %s", RemoteCatalogURL)
		remoteScripts, err := ListRemote()
		if err != nil {
			log.Printf("Error fetching remote catalog: %v", err)
			return scripts, fmt.Errorf("error fetching remote catalog: %w", err)
		}
		log.Printf("Found remote scripts: %v", remoteScripts)
		scripts = append(scripts, remoteScripts...)
	} else {
		log.Printf("No remote catalog URL set")
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
	if RemoteCatalogURL != "" {
		log.Printf("Fetching remote script: %s", name)
		bytes, err := ReadRemote(name)
		if err != nil {
			log.Printf("Error fetching remote script: %v", err)
			return nil, fmt.Errorf("script not found in embedded or remote catalog: %s", name)
		}
		return bytes, nil
	}

	return nil, fmt.Errorf("script not found: %s", name)
}
