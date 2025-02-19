package catalog

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Convert Go-style package reference to raw URL
func getRawURL(pkg string) (string, error) {
	// Basic validation
	if pkg == "" {
		return "", fmt.Errorf("empty package path")
	}

	parts := strings.Split(pkg, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid package path: %s", pkg)
	}

	// Is there a package that already does this?
	// Handle different Git hosting services
	switch parts[0] {
	case "github.com":
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/catalog", parts[1], parts[2]), nil
	case "gitlab.com":
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/raw/main/catalog", parts[1], parts[2]), nil
	case "bitbucket.org":
		return fmt.Sprintf("https://bitbucket.org/%s/%s/raw/main/catalog", parts[1], parts[2]), nil
	case "git.sr.ht": // SourceHut
		return fmt.Sprintf("https://git.sr.ht/%s/%s/blob/main/catalog", parts[1], parts[2]), nil
	case "codeberg.org":
		return fmt.Sprintf("https://codeberg.org/%s/%s/raw/branch/main/catalog", parts[1], parts[2]), nil
	default:
		return "", fmt.Errorf("unsupported repository host: %s", parts[0])
	}
}

// ListRemote fetches the list of available scripts from the remote catalog
func ListRemote() ([]string, error) {
	rawURL, err := getRawURL(viper.GetString("remote_catalog_url"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	fullURL := rawURL + "/index"

	// Use cached fetchURL instead of direct http.Get
	body, err := fetchURL(fullURL, cacheTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote catalog index: %w", err)
	}

	// Assuming the index is a newline-separated list of script names
	scripts := strings.Split(strings.TrimSpace(string(body)), "\n")
	return scripts, nil
}

// ReadRemote fetches a specific script from the remote catalog
func ReadRemote(name string) ([]byte, error) {
	// Sanitize the script name to prevent directory traversal
	name = filepath.Base(name)

	rawURL, err := getRawURL(viper.GetString("remote_catalog_url"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	// Use cached fetchURL instead of direct http.Get
	return fetchURL(rawURL+"/"+name, cacheTTL)
}
