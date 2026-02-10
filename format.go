package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func format(entry map[string]any) string {
	// Extract standard fields
	timeStr := formatTime(entry["time"])
	level := strings.ToUpper(fmt.Sprintf("%v", entry["level"]))
	msg := fmt.Sprintf("%v", entry["msg"])

	// Build attributes string
	var attrs []string
	keys := make([]string, 0, len(entry))
	for k := range entry {
		if k != "time" && k != "level" && k != "msg" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		attrs = append(attrs, fmt.Sprintf("%s=%v", k, entry[k]))
	}

	attrStr := ""
	if len(attrs) > 0 {
		attrStr = " " + strings.Join(attrs, " ")
	}

	return fmt.Sprintf("%s %s %s%s", timeStr, level, msg, attrStr)
}

func formatTime(t any) string {
	if t == nil {
		return ""
	}
	timeStr := fmt.Sprintf("%v", t)
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	return parsed.Format("2006-01-02 15:04:05")
}
