package catalog

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheDir = "run-catalog"
	cacheTTL = 10 * time.Minute
)

// CacheEntry represents a cached URL response
type CacheEntry struct {
	Expiration time.Time `json:"expiration"`
	Content    []byte    `json:"content"`
}

// getCacheDir returns the cache directory for run
func getCacheDir() (string, error) {
	// on macos, the cache directory is in ~/Library/Caches
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cachePath := filepath.Join(cacheDir, "run-catalog") // Using our constant
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", err
	}
	return cachePath, nil
}

func ClearCache() error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}
	fmt.Println("Clearing cache in", cacheDir)
	return os.RemoveAll(cacheDir)
}

// hashURL creates a hashed filename for a URL
func hashURL(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}

// getCachedResponse checks if a valid cached response exists
func getCachedResponse(url string) ([]byte, bool, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, false, err
	}

	cacheFile := filepath.Join(cacheDir, hashURL(url)+".json")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false, err
	}

	if time.Now().After(entry.Expiration) {
		_ = os.Remove(cacheFile) // Expired, remove it
		return nil, false, nil
	}

	return entry.Content, true, nil
}

// cacheResponse saves a response to the cache
func cacheResponse(url string, content []byte, ttl time.Duration) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(cacheDir, hashURL(url)+".json")

	entry := CacheEntry{
		Expiration: time.Now().Add(ttl),
		Content:    content,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0644)
}

// fetchURL fetches the URL, checking the cache first
func fetchURL(url string, ttl time.Duration) ([]byte, error) {
	if cachedData, found, err := getCachedResponse(url); err == nil && found {
		return cachedData, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status not OK: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := cacheResponse(url, body, ttl); err != nil {
		// Log cache persistence errors but don't fail the request
		fmt.Fprintf(os.Stderr, "Warning: failed to cache response: %v\n", err)
	}

	return body, nil
}
