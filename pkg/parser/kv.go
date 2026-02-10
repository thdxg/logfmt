package parser

import (
	"fmt"
	"strconv"
	"unicode"
)

// parseKV parses a line in key=value format.
// It requires at least one valid key=value pair to succeed.
// Tokens without '=' cause an error.
func parseKV(line []byte) (map[string]any, error) {
	raw := make(map[string]any)
	decoder := newKVDecoder(line)
	for decoder.scan() {
		raw[decoder.key] = decoder.val
	}
	if decoder.err != nil {
		return nil, decoder.err
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("no key=value pairs found")
	}
	return raw, nil
}

// kvDecoder is a simple key=value format decoder.
type kvDecoder struct {
	data    []byte
	pos     int
	key     string
	val     any
	err     error
	foundKV bool // tracks whether at least one key=value pair was found
}

func newKVDecoder(data []byte) *kvDecoder {
	return &kvDecoder{data: data}
}

func (d *kvDecoder) scan() bool {
	d.skipWhitespace()
	if d.pos >= len(d.data) {
		// If we reached the end without ever finding a key=value pair, report error
		if !d.foundKV {
			d.err = fmt.Errorf("no key=value pairs found")
		}
		return false
	}

	// Scan Key
	keyStart := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != '=' && d.data[d.pos] != ' ' {
		d.pos++
	}

	d.key = string(d.data[keyStart:d.pos])

	if d.pos >= len(d.data) || d.data[d.pos] == ' ' {
		// Token without '=' — this is not valid key=value format
		d.err = fmt.Errorf("unexpected token %q without '='", d.key)
		return false
	}

	// Found '='
	d.foundKV = true
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
func (d *kvDecoder) scanQuoted() string {
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

func (d *kvDecoder) skipWhitespace() {
	for d.pos < len(d.data) && unicode.IsSpace(rune(d.data[d.pos])) {
		d.pos++
	}
}
