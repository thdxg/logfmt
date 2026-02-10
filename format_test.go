package main

import (
	"testing"
)

type test struct {
	name  string
	entry map[string]any
	want  string
}

func Test_format(t *testing.T) {
	tests := []test{
		{
			name: "basic",
			entry: map[string]any{
				"time":  "2006-01-02T15:04:05.383759-05:00",
				"level": "INFO",
				"msg":   "Basic log",
			},
			want: "2006-01-02 15:04:05 INFO Basic log",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format(tt.entry)
			if got != tt.want {
				t.Errorf("\ngot  %vwant %v", got, tt.want)
			}
		})
	}
}
