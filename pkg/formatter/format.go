package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/thdxg/logfmt/pkg/color"
	"github.com/thdxg/logfmt/pkg/config"
	"github.com/thdxg/logfmt/pkg/types"
)

func Format(cfg config.Config, entry types.Entry) string {
	var timeStr string
	if !entry.Time.IsZero() {
		timeStr = entry.Time.Format(cfg.TimeFormat)
	} else {
		timeStr = entry.RawTime
	}

	level := formatLevel(entry.Level, cfg.LevelFormat, cfg.Color)
	msg := entry.Msg

	attrStr := ""
	if !cfg.HideAttrs {
		attrs := flattenMap(entry.Attrs)
		sort.Strings(attrs)
		if len(attrs) > 0 {
			rawAttrs := strings.Join(attrs, " ")
			if cfg.Color {
				attrStr = " " + color.Sprint(color.Gray, rawAttrs)
			} else {
				attrStr = " " + rawAttrs
			}
		}
	}

	return fmt.Sprintf("%s %s %s%s", timeStr, level, msg, attrStr)
}

func formatLevel(lvl string, style string, colorize bool) string {
	var output string
	switch style {
	case "short":
		switch lvl {
		case "INFO":
			output = "INF"
		case "WARN":
			output = "WRN"
		case "ERROR":
			output = "ERR"
		case "DEBUG":
			output = "DBG"
		default:
			if len(lvl) > 3 {
				output = lvl[:3]
			} else {
				output = lvl
			}
		}
	case "tiny":
		if len(lvl) > 0 {
			output = lvl[:1]
		}
	default:
		output = lvl
	}

	if !colorize {
		return output
	}

	switch lvl {
	case "INFO":
		return color.Sprint(color.Blue, output)
	case "WARN":
		return color.Sprint(color.Yellow, output)
	case "ERROR":
		return color.Sprint(color.Red, output)
	case "DEBUG":
		return color.Sprint(color.Green, output)
	}

	return output
}

func flattenMap(entry map[string]any) []string {
	var attrs []string

	var visit func(map[string]any, string)
	visit = func(m map[string]any, prefix string) {
		for k, v := range m {
			fullKey := k
			if prefix != "" {
				fullKey = prefix + "." + k
			}

			if nestedMap, ok := v.(map[string]any); ok {
				visit(nestedMap, fullKey)
			} else {
				attrs = append(attrs, fmt.Sprintf("%s=%v", fullKey, v))
			}
		}
	}

	visit(entry, "")

	return attrs
}
