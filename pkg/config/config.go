package config

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/thdxg/logfmt/pkg/types"
)

type Config struct {
	TimeFormat  string            `mapstructure:"time-format"`
	LevelFormat types.LevelFormat `mapstructure:"level-format"`
	Color       bool              `mapstructure:"color"`
	HideAttrs   bool              `mapstructure:"hide-attrs"`
}

func Default() Config {
	return Config{
		TimeFormat:  "2006-01-02 15:04:05",
		LevelFormat: types.LevelFormatFull,
		Color:       true,
		HideAttrs:   false,
	}
}

// Load reads config and returns struct.
// It sets up viper, binds flags (if provided), reads config file, and unmarshals.
func Load(cfgFile string, flags *pflag.FlagSet) (Config, error) {
	v := viper.New()

	// Defaults
	defaults := Default()
	v.SetDefault("time-format", defaults.TimeFormat)
	v.SetDefault("level-format", defaults.LevelFormat)
	v.SetDefault("color", defaults.Color)
	v.SetDefault("hide-attrs", defaults.HideAttrs)

	// Bind Flags
	if flags != nil {
		if err := v.BindPFlags(flags); err != nil {
			return defaults, err
		}
	}

	// Config File
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			v.AddConfigPath(home)
		}
		v.AddConfigPath(".")
		v.SetConfigType("yaml")
		v.SetConfigName(".logfmt")
	}

	// Env Vars
	v.SetEnvPrefix("LOGFMT")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Read Config (ignore errors if file not found)
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return defaults, err
	}

	// Validate LevelFormat
	if !cfg.LevelFormat.IsValid() {
		cfg.LevelFormat = defaults.LevelFormat
	}

	return cfg, nil
}
