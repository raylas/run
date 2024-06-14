package equip

import (
	"encoding/base64"
	"fmt"

	"github.com/linecard/run/catalog"
)

const (
	command = "echo %s | base64 -d -i > /run && chmod +x /run && /run %s"
)

func Pack(name string, args string) (string, error) {
	raw, err := catalog.Read(name)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(raw)

	return fmt.Sprintf(command, encoded, args), nil
}
