package cmd

import (
	"fmt"
	"log/slog"

	"github.com/USA-RedDragon/mandelbrot/internal/config"
	"github.com/USA-RedDragon/mandelbrot/internal/game"
	"github.com/USA-RedDragon/mandelbrot/internal/mandelbrot"
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
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	loader := config.New(cmd.Flags())
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		cfg, err := loader.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return run(cmd, cfg)
	}
	return cmd
}

func run(cmd *cobra.Command, cfg *config.Config) error {
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

	slog.Info("mandelbrot", "version", cmd.Annotations["version"], "commit", cmd.Annotations["commit"])

	palette := mandelbrot.PaletteModeSimpleRainbow
	if cfg.Palette == config.PaletteGrayscale {
		palette = mandelbrot.PaletteModeSimpleGrayscale
	}

	game, err := game.NewGame(cfg.Width, cfg.Height, mandelbrot.Settings{
		MaxIterations: cfg.MaxIterations,
		Scale:         cfg.Scale,
		Center:        cfg.Center,
		Exponent:      cfg.Exponent,
		StartingZ:     cfg.Z,
		StartingC:     cfg.C,
		Julia:         cfg.Julia,
		Palette:       palette,
	})
	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	if err := ebiten.RunGame(game); err != nil {
		return fmt.Errorf("failed to run game: %w", err)
	}

	return nil
}
