package formatter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/thdxg/logfmt/pkg/config"
)

func Format(entry map[string]any, cfg config.Config) string {
	// Extract standard fields
	timeStr := formatTime(entry["time"], cfg.TimeFormat)
	level := formatLevel(entry["level"], cfg.LevelFormat, cfg.Color)
	msg := fmt.Sprintf("%v", entry["msg"])

	// Build attributes string
	attrStr := ""
	if !cfg.HideAttrs {
		attrs := flattenMap(entry)
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

func formatTime(t any, layout string) string {
	if t == nil {
		return ""
	}
	timeStr := fmt.Sprintf("%v", t)
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return parsed.Format(layout)
}

func formatLevel(l any, style string, colorize bool) string {
	levelStr := strings.ToUpper(fmt.Sprintf("%v", l))

	var output string
	switch style {
	case "short":
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
	case "tiny":
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
			if k == "time" || k == "level" || k == "msg" {
				if prefix == "" {
					continue
				}
			}

			fullKey := k
			if prefix != "" {
				fullKey = prefix + "." + k
			}

			if nestedMap, ok := v.(map[string]any); ok {
				visit(nestedMap, fullKey)
			} else {
				// Handle arrays specifically if needed, or just standard formatting
				// User requested "Simple String" for arrays: key=[item1 item2]
				// Go's default %v does [item1 item2] for arrays, so that works.
				attrs = append(attrs, fmt.Sprintf("%s=%v", fullKey, v))
			}
		}
	}

	visit(entry, "")
	return attrs
}
