package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thdxg/logfmt/pkg/config"
	"github.com/thdxg/logfmt/pkg/formatter"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "logfmt",
	Short: "Formats structured logs from stdin",
	Long:  `logfmt reads structured logs (JSON) from stdin and formats them into a readable output.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(cfgFile, cmd.Flags())
		if err != nil {
			cobra.CheckErr(err)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var entry map[string]any
			if err := json.Unmarshal(line, &entry); err != nil {
				// If line is not valid JSON, print it as is
				fmt.Println(string(line))
				continue
			}
			fmt.Println(formatter.Format(entry, cfg))
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.logfmt.yaml)")

	// Define flags with defaults that match Config defaults
	defaults := config.Default()
	rootCmd.PersistentFlags().String("time-format", defaults.TimeFormat, "Timestamp format")
	rootCmd.PersistentFlags().String("level-format", defaults.LevelFormat, "Level format: full (INFO), short (INF), tiny (I)")
	rootCmd.PersistentFlags().Bool("color", defaults.Color, "Enable colored output")
	rootCmd.PersistentFlags().Bool("hide-attrs", defaults.HideAttrs, "Hide log attributes, show only time, level, and msg")
}
