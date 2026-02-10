package config

import (
	"reflect"
	"testing"

	"github.com/thdxg/logfmt/pkg/types"
)

func TestDefault(t *testing.T) {
	want := Config{
		TimeFormat:  "2006-01-02 15:04:05",
		LevelFormat: types.LevelFormatFull,
		Color:       true,
		HideAttrs:   false,
	}
	if got := Default(); !reflect.DeepEqual(got, want) {
		t.Errorf("Default() = %v, want %v", got, want)
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Test loading with no file and no flags (should return defaults)
	cfg, err := Load("", nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := Default()
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %v, want %v", cfg, want)
	}
}
