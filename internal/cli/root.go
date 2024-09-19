package cli

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/pomdtr/sunbeam/internal/tui"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

func IsSunbeamRunning() bool {
	return len(os.Getenv("SUNBEAM")) > 0
}

func NewRootCmd() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:          "sunbeam",
		Short:        "Command Line Launcher",
		SilenceUsage: true,
		Long: `Sunbeam is a command line launcher for your terminal, inspired by fzf and raycast.

See https://pomdtr.github.io/sunbeam for more information.`,
	}

	rootCmd.AddGroup(&cobra.Group{
		ID:    "extension",
		Title: "Extension Commands:",
	})

	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if !strings.HasPrefix(entry.Name(), "sunbeam-") {
				continue
			}

			execPath := filepath.Join(dir, entry.Name())

			rootCmd.AddCommand(&cobra.Command{
				Use:                strings.TrimPrefix(entry.Name(), "sunbeam-"),
				DisableFlagParsing: true,
				Short:              strings.Replace(execPath, os.Getenv("HOME"), "~", 1),
				GroupID:            "extension",
				RunE: func(cmd *cobra.Command, args []string) error {
					runner := tui.NewRunner(execPath, args...)
					return tui.Draw(runner)
				},
			})
		}
	}

	return rootCmd, nil
}
