package catalog

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
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
	rawURL, err := getRawURL(RemoteCatalogURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	fullURL := rawURL + "/index"
	log.Printf("Fetching remote index from: %s", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP status not OK: %s", resp.Status)
		return nil, fmt.Errorf("failed to fetch remote catalog index: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	content := string(body)
	log.Printf("Raw index content: %q", content)

	// Assuming the index is a newline-separated list of script names
	scripts := strings.Split(strings.TrimSpace(content), "\n")
	log.Printf("Parsed scripts: %v", scripts)

	return scripts, nil
}

// ReadRemote fetches a specific script from the remote catalog
func ReadRemote(name string) ([]byte, error) {
	// Sanitize the script name to prevent directory traversal
	name = filepath.Base(name)

	rawURL, err := getRawURL(RemoteCatalogURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	resp, err := http.Get(rawURL + "/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch remote script: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
