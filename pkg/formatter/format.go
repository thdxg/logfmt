package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/thdxg/logfmt/pkg/config"
	"github.com/thdxg/logfmt/pkg/types"
)

func Format(entry types.Entry, cfg config.Config) string {
	// Extract standard fields
	var timeStr string
	if !entry.Time.IsZero() {
		// Use standard time format if parsed successfully
		timeStr = entry.Time.Format(cfg.TimeFormat)
	} else {
		// Fallback to raw time string if parsing failed
		timeStr = entry.RawTime
	}

	level := formatLevel(entry.Level, cfg.LevelFormat, cfg.Color)
	msg := entry.Msg // Already a string

	// Build attributes string
	attrStr := ""
	if !cfg.HideAttrs {
		attrs := flattenMap(entry.Attrs)
		sort.Strings(attrs)
		if len(attrs) > 0 {
			rawAttrs := strings.Join(attrs, " ")
			if cfg.Color {
				attrStr = " " + color.New(color.FgHiBlack).Sprint(rawAttrs)
			} else {
				attrStr = " " + rawAttrs
			}
		}
	}

	return fmt.Sprintf("%s %s %s%s", timeStr, level, msg, attrStr)
}

func formatLevel(lvl types.Level, style types.LevelFormat, colorize bool) string {
	levelStr := strings.ToUpper(string(lvl))

	var output string
	switch style {
	case types.LevelFormatShort:
		switch levelStr {
		case "INFO":
			output = "INF"
		case "WARN", "WARNING":
			output = "WRN"
		case "ERROR":
			output = "ERR"
		case "DEBUG":
			output = "DBG"
		default:
			if len(levelStr) > 3 {
				output = levelStr[:3]
			} else {
				output = levelStr
			}
		}
	case types.LevelFormatTiny:
		if len(levelStr) > 0 {
			output = levelStr[:1]
		}
	default: // full
		output = levelStr
	}

	if !colorize {
		return output
	}

	c := color.New()
	switch levelStr {
	case "INFO":
		c.Add(color.FgBlue)
	case "WARN", "WARNING":
		c.Add(color.FgYellow)
	case "ERROR":
		c.Add(color.FgRed)
	case "DEBUG":
		c.Add(color.FgGreen)
	}
	return c.Sprint(output)
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
