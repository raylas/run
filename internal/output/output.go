package output

import "fmt"

type Color string

const (
	Reset   Color = "\x1b[0000m"
	Default       = "\x1b[0039m"
	Green         = "\x1b[0032m"
	Red           = "\x1b[0031m"
)

func Colorize(color Color, format string, a ...interface{}) string {
	return fmt.Sprintf("%v%v%v", color, fmt.Sprintf(format, a...), Reset)
}
