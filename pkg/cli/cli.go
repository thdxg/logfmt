package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/thdxg/logfmt/pkg/config"
	"github.com/thdxg/logfmt/pkg/types"
)

// ParseFlags parses command line flags and returns a config.Config
// with the correct precedence: Flags > Env > Default.
func ParseFlags() config.Config {
	defaults := config.Default()

	flag.String("time-format", "", "Timestamp format (default: "+defaults.TimeFormat+")")
	flag.String("level-format", "", "Level format: full, short, tiny (default: "+string(defaults.LevelFormat)+")")
	flag.Bool("no-color", false, "Disable colored output")
	flag.Bool("hide-attrs", false, "Hide log attributes, show only time, level, and msg")

	flag.Usage = Usage
	flag.Parse()

	var timeFormat *string
	var levelFormat *string
	var colorPtr *bool
	var hideAttrsPtr *bool

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "time-format":
			v := f.Value.String()
			timeFormat = &v
		case "level-format":
			v := f.Value.String()
			levelFormat = &v
		case "no-color":
			v := false
			colorPtr = &v
		case "hide-attrs":
			v := true
			hideAttrsPtr = &v
		}
	})

	return config.Load(timeFormat, levelFormat, colorPtr, hideAttrsPtr)
}

func HasStdin() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func Usage() {
	defaults := config.Default()

	fmt.Fprintln(os.Stderr, "logfmt - formats structured logs (JSON/key=value) from stdin")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage: <command> | logfmt [flags]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	fmt.Fprintf(os.Stderr, "  -time-format string\n")
	fmt.Fprintf(os.Stderr, "    \tTimestamp format (default: %s)\n", defaults.TimeFormat)
	fmt.Fprintf(os.Stderr, "  -level-format string\n")
	fmt.Fprintf(os.Stderr, "    \tLevel format: %s, %s, %s (default: %s)\n",
		types.LevelFormatFull, types.LevelFormatShort, types.LevelFormatTiny, defaults.LevelFormat)
	fmt.Fprintf(os.Stderr, "  -no-color\n")
	fmt.Fprintf(os.Stderr, "    \tDisable colored output\n")
	fmt.Fprintf(os.Stderr, "  -hide-attrs\n")
	fmt.Fprintf(os.Stderr, "    \tHide log attributes, show only time, level, and msg\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Environment variables:")
	fmt.Fprintln(os.Stderr, "  LOGFMT_TIME_FORMAT    Timestamp format")
	fmt.Fprintln(os.Stderr, "  LOGFMT_LEVEL_FORMAT   Level format: full, short, tiny")
	fmt.Fprintln(os.Stderr, "  LOGFMT_COLOR          Enable/disable colored output (true/false)")
	fmt.Fprintln(os.Stderr, "  LOGFMT_HIDE_ATTRS     Hide log attributes (true/false)")
}
