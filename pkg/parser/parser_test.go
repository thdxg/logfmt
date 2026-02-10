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
			name:  "json with nested attributes",
			input: `{"time":"2023-01-01T10:00:00Z", "level":"INFO", "msg":"nested", "user":{"name":"Alice"}}`,
			check: func(t *testing.T, e types.Entry) {
				if e.Msg != "nested" {
					t.Errorf("expected msg 'nested', got %v", e.Msg)
				}
				user, ok := e.Attrs["user"].(map[string]any)
				if !ok {
					t.Fatal("expected user to be a map")
				}
				if user["name"] != "Alice" {
					t.Errorf("expected user.name 'Alice', got %v", user["name"])
				}
			},
		},
		{
			name:  "kv simple",
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
			name:  "kv quoted values",
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
			name:    "plain text is rejected",
			input:   "plain text message",
			wantErr: true,
		},
		{
			name:    "empty line is rejected",
			input:   "",
			wantErr: true,
		},
		{
			name:    "bare token in kv is rejected",
			input:   "level=INFO bare msg=hello",
			wantErr: true,
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
