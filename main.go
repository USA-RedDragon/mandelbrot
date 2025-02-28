package main

import (
	"log/slog"
	"os"

	"github.com/USA-RedDragon/mandelbrot/internal/cmd"
)

// https://goreleaser.com/cookbooks/using-main.version/
//
//nolint:golint,gochecknoglobals
var (
	version = "dev"
	commit  = "none"
)

func main() {
	rootCmd := cmd.NewCommand(version, commit)
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute command", "error", err)
		os.Exit(1)
	}
}
