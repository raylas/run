package equip

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/viper"
)

func Pack(script []byte, args string) string {
	encoded := base64.StdEncoding.EncodeToString(script)
	return fmt.Sprintf(viper.GetString("command"), encoded, args)
}
