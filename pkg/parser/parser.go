package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/thdxg/logfmt/pkg/types"
)

// Parse attempts to parse a log line into an Entry.
// It detects JSON vs logfmt based on the first character.
func Parse(line []byte) (types.Entry, error) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return types.Entry{}, fmt.Errorf("empty line")
	}

	if line[0] == '{' {
		// Try JSON
		var raw map[string]any
		if err := json.Unmarshal(line, &raw); err == nil {
			return toEntry(raw), nil
		}
		// If JSON fails, fallback to logfmt (maybe it's text starting with {?)
	}
	return parseLogfmt(line)
}

func parseLogfmt(line []byte) (types.Entry, error) {
	raw := make(map[string]any)
	decoder := newLogfmtDecoder(line)
	for decoder.scan() {
		raw[decoder.key] = decoder.val
	}
	if decoder.err != nil {
		return types.Entry{}, decoder.err
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
			// Try parsing time
			if t, err := parseTime(v); err == nil {
				e.Time = t
			}
		case "level", "lvl", "severity":
			e.Level = types.Level(fmt.Sprintf("%v", v))
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
	// Try standard layouts
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

// Simple logfmt decoder implementation
type logfmtDecoder struct {
	data []byte
	pos  int
	key  string
	val  any
	err  error
}

func newLogfmtDecoder(data []byte) *logfmtDecoder {
	return &logfmtDecoder{data: data}
}

func (d *logfmtDecoder) scan() bool {
	d.skipWhitespace()
	if d.pos >= len(d.data) {
		return false
	}

	// Scan Key
	keyStart := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != '=' && d.data[d.pos] != ' ' {
		d.pos++
	}

	d.key = string(d.data[keyStart:d.pos])

	if d.pos >= len(d.data) || d.data[d.pos] == ' ' {
		// Key without value (boolean flag)
		d.val = true
		return true
	}

	// Found '='
	d.pos++ // skip '='

	// Scan Value
	if d.pos >= len(d.data) {
		d.val = ""
		return true
	}

	if d.data[d.pos] == '"' {
		// Quoted value
		quotedVal, err := strconv.Unquote(d.scanQuoted())
		if err != nil {
			d.err = err
			return false
		}
		d.val = quotedVal
	} else {
		// Unquoted value
		valStart := d.pos
		for d.pos < len(d.data) && d.data[d.pos] != ' ' {
			d.pos++
		}
		d.val = string(d.data[valStart:d.pos])
	}
	return true
}

// scanQuoted returns the raw quoted string including quotes
func (d *logfmtDecoder) scanQuoted() string {
	start := d.pos
	d.pos++ // skip opening quote
	escaped := false
	for d.pos < len(d.data) {
		c := d.data[d.pos]
		if escaped {
			escaped = false
		} else if c == '\\' {
			escaped = true
		} else if c == '"' {
			d.pos++
			return string(d.data[start:d.pos])
		}
		d.pos++
	}
	return string(d.data[start:]) // Return potentially incomplete string to let Unquote handle error
}

func (d *logfmtDecoder) skipWhitespace() {
	for d.pos < len(d.data) && unicode.IsSpace(rune(d.data[d.pos])) {
		d.pos++
	}
}
