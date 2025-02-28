package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel LogLevel `json:"log-level" yaml:"log-level"`
	Width    uint     `json:"width" yaml:"width"`
	Height   uint     `json:"height" yaml:"height"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

//nolint:golint,gochecknoglobals
var (
	ConfigFileKey = "config"
	LogLevelKey   = "log-level"
	WidthKey      = "width"
	HeightKey     = "height"
)

const (
	DefaultConfigPath = "config.yaml"
	DefaultLogLevel   = LogLevelInfo
	DefaultWidth      = 720
	DefaultHeight     = 480
)

func RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(ConfigFileKey, "c", DefaultConfigPath, "Config file path")
	cmd.Flags().String(LogLevelKey, string(DefaultLogLevel), "Log level")
	cmd.Flags().Uint(WidthKey, DefaultWidth, "Initial window width")
	cmd.Flags().Uint(HeightKey, DefaultHeight, "Initial window height")
}

var (
	ErrInvalidLogLevel = errors.New("Invalid log level")
	ErrInvalidWidth    = errors.New("Invalid width")
	ErrInvalidHeight   = errors.New("Invalid height")
)

func (c *Config) Validate() error {
	switch c.LogLevel {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
	default:
		return ErrInvalidLogLevel
	}

	if c.Width == 0 {
		return ErrInvalidWidth
	}

	if c.Height == 0 {
		return ErrInvalidHeight
	}

	return nil
}

func LoadConfig(cmd *cobra.Command) (*Config, error) {
	var config Config

	// Load flags from envs
	ctx, cancel := context.WithCancelCause(cmd.Context())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if ctx.Err() != nil {
			return
		}
		optName := strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_"), ".", "__")
		if val, ok := os.LookupEnv(optName); !f.Changed && ok {
			if err := f.Value.Set(val); err != nil {
				cancel(err)
			}
			f.Changed = true
		}
	})
	if ctx.Err() != nil {
		return &config, fmt.Errorf("failed to load env: %w", context.Cause(ctx))
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return &config, fmt.Errorf("failed to get config path: %w", err)
	}
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return &config, fmt.Errorf("failed to read config: %w", err)
		} else if err == nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				return &config, fmt.Errorf("failed to unmarshal config: %w", err)
			}
		}
	}

	err = overrideFlags(&config, cmd)
	if err != nil {
		return &config, fmt.Errorf("failed to override flags: %w", err)
	}

	// Defaults
	if config.LogLevel == "" {
		config.LogLevel = DefaultLogLevel
	}

	if config.Width == 0 {
		config.Width = DefaultWidth
	}

	if config.Height == 0 {
		config.Height = DefaultHeight
	}

	return &config, nil
}

func overrideFlags(config *Config, cmd *cobra.Command) error {
	if cmd.Flags().Changed(LogLevelKey) {
		ll, err := cmd.Flags().GetString(LogLevelKey)
		if err != nil {
			return fmt.Errorf("failed to get log level: %w", err)
		}
		config.LogLevel = LogLevel(ll)
	}

	if cmd.Flags().Changed(WidthKey) {
		w, err := cmd.Flags().GetUint(WidthKey)
		if err != nil {
			return fmt.Errorf("failed to get width: %w", err)
		}
		config.Width = w
	}

	if cmd.Flags().Changed(HeightKey) {
		h, err := cmd.Flags().GetUint(HeightKey)
		if err != nil {
			return fmt.Errorf("failed to get height: %w", err)
		}
		config.Height = h
	}

	return nil
}
