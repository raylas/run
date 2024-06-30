package parse

import (
	"bufio"
	"regexp"
	"strings"
)

type Arg struct {
	Name  string
	Value string
	Desc  string
}

func Args(script []byte) (string, map[int]Arg, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(script)))

	rd := regexp.MustCompile(`(?m)^## (.*)$`)
	ra := regexp.MustCompile(`(?m)^# (\w+): (.+?)(?: \[(.+?)\])?$`)

	desc := ""
	pos := 0
	args := make(map[int]Arg)

	for scanner.Scan() {
		line := scanner.Text()

		descMatches := rd.FindStringSubmatch(line)
		if len(descMatches) == 2 {
			desc = string(descMatches[1])
			continue
		}

		argMatches := ra.FindStringSubmatch(line)
		if len(argMatches) > 0 {
			args[pos] = Arg{
				Name:  strings.ToLower(argMatches[1]),
				Value: argMatches[3],
				Desc:  argMatches[2],
			}
			pos++
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, err
	}

	return desc, args, nil
}
