package config

import (
	"os"
	"strconv"

	"github.com/thdxg/logfmt/pkg/types"
)

type Config struct {
	TimeFormat  string
	LevelFormat string
	Color       bool
	HideAttrs   bool
}

func Default() Config {
	return Config{
		TimeFormat:  "2006-01-02 15:04:05",
		LevelFormat: "full",
		Color:       true,
		HideAttrs:   false,
	}
}

// Load loads configuration with precedence: Flags > Env > Default.
// Nil pointers indicate the value was not set by the caller.
func Load(timeFormat *string, levelFormat *string, color *bool, hideAttrs *bool) Config {
	cfg := Default()

	// 1. Env Vars
	if v := os.Getenv("LOGFMT_TIME_FORMAT"); v != "" {
		cfg.TimeFormat = v
	}
	if v := os.Getenv("LOGFMT_LEVEL_FORMAT"); v != "" {
		cfg.LevelFormat = v
	}
	if v := os.Getenv("LOGFMT_COLOR"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Color = b
		}
	}
	if v := os.Getenv("LOGFMT_HIDE_ATTRS"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.HideAttrs = b
		}
	}

	// 2. Flags (override env if explicitly set)
	if timeFormat != nil {
		cfg.TimeFormat = *timeFormat
	}
	if levelFormat != nil {
		cfg.LevelFormat = *levelFormat
	}
	if color != nil {
		cfg.Color = *color
	}
	if hideAttrs != nil {
		cfg.HideAttrs = *hideAttrs
	}

	// Validation
	switch cfg.LevelFormat {
	case types.LevelFormatFull:
	case types.LevelFormatShort:
	case types.LevelFormatTiny:
	default:
		cfg.LevelFormat = types.LevelFormatFull
	}

	return cfg
}
