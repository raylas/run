package catalog

import (
	"embed"
	"fmt"
)

//go:embed *
var catalog embed.FS

func List() ([]string, error) {
	files, err := catalog.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var scripts []string
	for _, file := range files {
		if file.Name() == "catalog.go" || file.IsDir() {
			continue
		}

		scripts = append(scripts, file.Name())
	}

	return scripts, nil
}

func Read(name string) ([]byte, error) {
	bytes, err := catalog.ReadFile(name)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, fmt.Errorf("script has no bytes: %s", name)
	}

	return bytes, nil
}
