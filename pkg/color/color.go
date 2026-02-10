package color

import (
	"fmt"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Gray   = "\033[90m"
)

// Sprint returns the string wrapped in the given color code.
// It unconditionally applies the color code. The caller is responsible for checking if color should be enabled.
func Sprint(code string, msg any) string {
	return fmt.Sprintf("%s%v%s", code, msg, Reset)
}
