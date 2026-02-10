package types

import "time"

type Entry struct {
	Time    time.Time
	RawTime string
	Level   string
	Msg     string
	Attrs   map[string]any
}
