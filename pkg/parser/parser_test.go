package parser

import (
	"testing"

	"github.com/thdxg/logfmt/pkg/types"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, types.Entry)
	}{
		{
			name:  "json simple",
			input: `{"time":"2023-01-01T10:00:00Z", "level":"INFO", "msg":"hello"}`,
			check: func(t *testing.T, e types.Entry) {
				if e.Msg != "hello" {
					t.Errorf("expected msg 'hello', got %v", e.Msg)
				}
				if string(e.Level) != "INFO" {
					t.Errorf("expected level 'INFO', got %v", e.Level)
				}
				if e.Time.IsZero() {
					t.Error("expected time to be parsed")
				}
			},
		},
		{
			name:  "logfmt simple",
			input: `time=2023-01-01T10:00:00Z level=INFO msg=hello`,
			check: func(t *testing.T, e types.Entry) {
				if e.Msg != "hello" {
					t.Errorf("expected msg 'hello', got %v", e.Msg)
				}
				if string(e.Level) != "INFO" {
					t.Errorf("expected level 'INFO', got %v", e.Level)
				}
				if e.Time.IsZero() {
					t.Error("expected time to be parsed")
				}
			},
		},
		{
			name:  "invalid time fallback",
			input: `time="invalid-time" level=INFO msg=test`,
			check: func(t *testing.T, e types.Entry) {
				if !e.Time.IsZero() {
					t.Error("expected zero time for invalid input")
				}
				if e.RawTime != "invalid-time" {
					t.Errorf("expected RawTime 'invalid-time', got %v", e.RawTime)
				}
			},
		},
		{
			name:  "logfmt quoted",
			input: `level=WARN msg="hello world" key="value with spaces"`,
			check: func(t *testing.T, e types.Entry) {
				if e.Msg != "hello world" {
					t.Errorf("expected msg 'hello world', got %v", e.Msg)
				}
				if string(e.Level) != "WARN" {
					t.Errorf("expected level 'WARN', got %v", e.Level)
				}
				if e.Attrs["key"] != "value with spaces" {
					t.Errorf("expected attr key='value with spaces', got %v", e.Attrs["key"])
				}
			},
		},
		{
			name:  "logfmt boolean flag",
			input: `level=DEBUG msg=test debug_mode`,
			check: func(t *testing.T, e types.Entry) {
				if e.Msg != "test" {
					t.Errorf("expected msg 'test', got %v", e.Msg)
				}
				if e.Attrs["debug_mode"] != true {
					t.Errorf("expected attr debug_mode=true, got %v", e.Attrs["debug_mode"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
