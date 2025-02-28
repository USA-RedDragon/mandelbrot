package cmd

import (
	"fmt"
	"log/slog"

	"github.com/USA-RedDragon/mandelbrot/internal/config"
	"github.com/USA-RedDragon/mandelbrot/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func NewCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mandelbrot",
		Version: fmt.Sprintf("%s - %s", version, commit),
		Annotations: map[string]string{
			"version": version,
			"commit":  commit,
		},
		RunE:          run,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	config.RegisterFlags(cmd)
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	slog.Info("mandelbrot", "version", cmd.Annotations["version"], "commit", cmd.Annotations["commit"])

	cfg, err := config.LoadConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch cfg.LogLevel {
	case config.LogLevelDebug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case config.LogLevelInfo:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case config.LogLevelWarn:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case config.LogLevelError:
		slog.SetLogLoggerLevel(slog.LevelError)
	}

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	game, err := game.NewGame(720, 480)
	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	if err := ebiten.RunGame(game); err != nil {
		return fmt.Errorf("failed to run game: %w", err)
	}

	return nil
}
