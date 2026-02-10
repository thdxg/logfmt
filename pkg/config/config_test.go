package config

import (
	"testing"

	"github.com/thdxg/logfmt/pkg/types"
)

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func TestDefault(t *testing.T) {
	want := Config{
		TimeFormat:  "2006-01-02 15:04:05",
		LevelFormat: types.LevelFormatFull,
		Color:       true,
		HideAttrs:   false,
	}
	got := Default()
	if got != want {
		t.Errorf("Default() = %v, want %v", got, want)
	}
}

func TestLoad_Defaults(t *testing.T) {
	cfg := Load(nil, nil, nil, nil)
	want := Default()
	if cfg != want {
		t.Errorf("Load() = %v, want %v", cfg, want)
	}
}

func TestLoad_Flags(t *testing.T) {
	cfg := Load(strPtr("15:04"), strPtr("short"), boolPtr(false), boolPtr(true))

	if cfg.TimeFormat != "15:04" {
		t.Errorf("TimeFormat = %q, want %q", cfg.TimeFormat, "15:04")
	}
	if cfg.LevelFormat != types.LevelFormatShort {
		t.Errorf("LevelFormat = %q, want %q", cfg.LevelFormat, types.LevelFormatShort)
	}
	if cfg.Color != false {
		t.Errorf("Color = %v, want false", cfg.Color)
	}
	if cfg.HideAttrs != true {
		t.Errorf("HideAttrs = %v, want true", cfg.HideAttrs)
	}
}

func TestLoad_EnvVars(t *testing.T) {
	t.Setenv("LOGFMT_TIME_FORMAT", "15:04:05")
	t.Setenv("LOGFMT_LEVEL_FORMAT", "tiny")
	t.Setenv("LOGFMT_COLOR", "false")
	t.Setenv("LOGFMT_HIDE_ATTRS", "true")

	cfg := Load(nil, nil, nil, nil)

	if cfg.TimeFormat != "15:04:05" {
		t.Errorf("TimeFormat = %q, want %q", cfg.TimeFormat, "15:04:05")
	}
	if cfg.LevelFormat != types.LevelFormatTiny {
		t.Errorf("LevelFormat = %q, want %q", cfg.LevelFormat, types.LevelFormatTiny)
	}
	if cfg.Color != false {
		t.Errorf("Color = %v, want false", cfg.Color)
	}
	if cfg.HideAttrs != true {
		t.Errorf("HideAttrs = %v, want true", cfg.HideAttrs)
	}
}

func TestLoad_FlagsOverrideEnv(t *testing.T) {
	t.Setenv("LOGFMT_LEVEL_FORMAT", "tiny")

	cfg := Load(nil, strPtr("short"), nil, nil)

	if cfg.LevelFormat != types.LevelFormatShort {
		t.Errorf("LevelFormat = %q, want %q (flag should override env)", cfg.LevelFormat, types.LevelFormatShort)
	}
}

func TestLoad_InvalidLevelFormat(t *testing.T) {
	cfg := Load(nil, strPtr("invalid"), nil, nil)

	if cfg.LevelFormat != types.LevelFormatFull {
		t.Errorf("LevelFormat = %q, want %q (invalid should fallback to default)", cfg.LevelFormat, types.LevelFormatFull)
	}
}
