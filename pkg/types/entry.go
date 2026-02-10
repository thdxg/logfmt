package types

import "time"

type Entry struct {
	Time    time.Time
	RawTime string
	Level   Level
	Msg     string
	Attrs   map[string]any
}
