package formatter

import (
	"testing"
	"time"

	"github.com/thdxg/logfmt/pkg/config"
	"github.com/thdxg/logfmt/pkg/types"
)

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

type test struct {
	name   string
	entry  types.Entry
	config config.Config
	want   string
}

func Test_Format(t *testing.T) {
	tests := []test{
		{
			name: "basic",
			entry: types.Entry{
				Time:  mustParseTime(time.RFC3339, "2006-01-02T15:04:05Z"),
				Level: "INFO",
				Msg:   "Basic log",
			},
			config: config.Config{TimeFormat: "2006-01-02 15:04:05", LevelFormat: "full"},
			want:   "2006-01-02 15:04:05 INFO Basic log",
		},
		{
			name: "short level",
			entry: types.Entry{
				Time:  mustParseTime(time.RFC3339, "2006-01-02T15:04:05Z"),
				Level: "INFO",
				Msg:   "Short level",
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: types.LevelFormatShort},
			want:   "15:04 INF Short level",
		},
		{
			name: "tiny level",
			entry: types.Entry{
				Time:  mustParseTime(time.RFC3339, "2006-01-02T15:04:05Z"),
				Level: "WARNING",
				Msg:   "Tiny level",
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: types.LevelFormatTiny},
			want:   "15:04 W Tiny level",
		},
		{
			name: "nested attributes",
			entry: types.Entry{
				Time:  mustParseTime(time.RFC3339, "2006-01-02T15:04:05Z"),
				Level: "INFO",
				Msg:   "Nested",
				Attrs: map[string]any{
					"user": map[string]any{
						"name": "Alice",
						"id":   123,
					},
				},
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: types.LevelFormatFull},
			want:   "15:04 INFO Nested user.id=123 user.name=Alice",
		},
		{
			name: "hide attributes",
			entry: types.Entry{
				Time:  mustParseTime(time.RFC3339, "2006-01-02T15:04:05Z"),
				Level: "INFO",
				Msg:   "Hidden",
				Attrs: map[string]any{
					"foo": "bar",
				},
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: types.LevelFormatFull, HideAttrs: true},
			want:   "15:04 INFO Hidden",
		},
		{
			name: "array attributes",
			entry: types.Entry{
				Time:  mustParseTime(time.RFC3339, "2006-01-02T15:04:05Z"),
				Level: "INFO",
				Msg:   "Array",
				Attrs: map[string]any{
					"tags": []any{"a", "b"},
				},
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: types.LevelFormatFull},
			want:   "15:04 INFO Array tags=[a b]",
		},
		{
			name: "raw time fallback",
			entry: types.Entry{
				Time:    time.Time{},
				RawTime: "invalid-time",
				Level:   "INFO",
				Msg:     "Fallback",
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: types.LevelFormatFull},
			want:   "invalid-time INFO Fallback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.config, tt.entry)
			if got != tt.want {
				t.Errorf("\ngot  %q\nwant %q", got, tt.want)
			}
		})
	}
}
