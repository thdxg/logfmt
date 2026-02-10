package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/thdxg/logfmt/pkg/types"
)

// Parse attempts to parse a log line into an Entry.
// It tries JSON first (if line starts with '{'), then key=value format.
// Returns an error if the line cannot be parsed as either format.
func Parse(line []byte) (types.Entry, error) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return types.Entry{}, fmt.Errorf("empty line")
	}

	if line[0] == '{' {
		var raw map[string]any
		if err := json.Unmarshal(line, &raw); err == nil {
			return toEntry(raw), nil
		}
	}

	raw, err := parseKV(line)
	if err != nil {
		return types.Entry{}, err
	}

	return toEntry(raw), nil
}

func toEntry(raw map[string]any) types.Entry {
	e := types.Entry{
		Attrs: make(map[string]any),
	}

	for k, v := range raw {
		switch strings.ToLower(k) {
		case "time", "t", "timestamp", "date":
			e.RawTime = fmt.Sprintf("%v", v)
			if t, err := parseTime(v); err == nil {
				e.Time = t
			}
		case "level", "lvl", "severity":
			e.Level = fmt.Sprintf("%v", v)
		case "msg", "message", "v", "val":
			e.Msg = fmt.Sprintf("%v", v)
		default:
			e.Attrs[k] = v
		}
	}
	return e
}

func parseTime(v any) (time.Time, error) {
	s := fmt.Sprintf("%v", v)
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown time format")
}
