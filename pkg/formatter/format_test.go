package formatter

import (
	"testing"

	"github.com/thdxg/logfmt/pkg/config"
)

type test struct {
	name   string
	entry  map[string]any
	config config.Config
	want   string
}

func Test_Format(t *testing.T) {
	tests := []test{
		{
			name: "basic",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05.383759-05:00",
				"level": "INFO",
				"msg":   "Basic log",
			},
			config: config.Config{TimeFormat: "2006-01-02 15:04:05", LevelFormat: "full"},
			want:   "2006-01-02 15:04:05 INFO Basic log",
		},
		{
			name: "short level",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05Z",
				"level": "INFO",
				"msg":   "Short level",
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: "short"},
			want:   "15:04 INF Short level",
		},
		{
			name: "tiny level",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05Z",
				"level": "WARNING",
				"msg":   "Tiny level",
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: "tiny"},
			want:   "15:04 W Tiny level",
		},
		{
			name: "nested attributes",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05Z",
				"level": "INFO",
				"msg":   "Nested",
				"user": map[string]any{
					"name": "Alice",
					"id":   123,
				},
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: "full"},
			want:   "15:04 INFO Nested user.id=123 user.name=Alice",
		},
		{
			name: "hide attributes",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05Z",
				"level": "INFO",
				"msg":   "Hidden",
				"foo":   "bar",
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: "full", HideAttrs: true},
			want:   "15:04 INFO Hidden",
		},
		{
			name: "array attributes",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05Z",
				"level": "INFO",
				"msg":   "Array",
				"tags":  []any{"a", "b"},
			},
			config: config.Config{TimeFormat: "15:04", LevelFormat: "full"},
			want:   "15:04 INFO Array tags=[a b]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.entry, tt.config)
			if got != tt.want {
				t.Errorf("\ngot  %q\nwant %q", got, tt.want)
			}
		})
	}
}
