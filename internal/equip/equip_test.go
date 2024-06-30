package equip

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	command := "echo %s | base64 -d -i > /run && chmod +x /run && /run %s"
	viper.Set("command", command)

	cases := []struct {
		desc     string
		script   []byte
		args     string
		expected string
	}{
		{
			desc:     "simple script",
			script:   []byte("echo 'Hello, world!'"),
			args:     "arg1 arg2",
			expected: fmt.Sprintf(command, "ZWNobyAnSGVsbG8sIHdvcmxkISc=", "arg1 arg2"),
		},
		{
			desc:     "empty script and args",
			script:   []byte(""),
			args:     "",
			expected: fmt.Sprintf(command, "", ""),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			packed := Pack(tc.script, tc.args)

			assert.Equal(t, tc.expected, packed)
		})
	}
}
