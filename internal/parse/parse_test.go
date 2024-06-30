package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgs(t *testing.T) {
	cases := []struct {
		desc         string
		script       []byte
		expectedDesc string
		expectedArgs map[int]Arg
		expectedErr  error
	}{
		{
			desc:         "empty script",
			script:       []byte(""),
			expectedDesc: "",
			expectedArgs: map[int]Arg{},
			expectedErr:  nil,
		},
		{
			desc:         "only description",
			script:       []byte("## This is a script description"),
			expectedDesc: "This is a script description",
			expectedArgs: map[int]Arg{},
			expectedErr:  nil,
		},
		{
			desc:         "multiple arguments",
			script:       []byte("## Script\n# arg1: Description 1 [default1]\n# arg2: Description 2 [default2]"),
			expectedDesc: "Script",
			expectedArgs: map[int]Arg{
				0: {Name: "arg1", Value: "default1", Desc: "Description 1"},
				1: {Name: "arg2", Value: "default2", Desc: "Description 2"},
			},
			expectedErr: nil,
		},
		{
			desc:   "only arguments",
			script: []byte("# arg1: Description 1 [default1]\n# arg2: Description 2 [default2]"),
			expectedArgs: map[int]Arg{
				0: {Name: "arg1", Value: "default1", Desc: "Description 1"},
				1: {Name: "arg2", Value: "default2", Desc: "Description 2"},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			desc, args, err := Args(tc.script)

			assert.Equal(t, tc.expectedDesc, desc)
			assert.Equal(t, tc.expectedArgs, args)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
