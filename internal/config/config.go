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
	LogLevelKey   = "log_level"
)

const (
	DefaultConfigPath = "config.yaml"
	DefaultLogLevel   = LogLevelInfo
)

func RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(ConfigFileKey, "c", DefaultConfigPath, "Config file path")
	cmd.Flags().String(LogLevelKey, string(DefaultLogLevel), "Log level")
}

var (
	ErrInvalidLogLevel = errors.New("Invalid log level")
)

func (c *Config) Validate() error {
	switch c.LogLevel {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
	default:
		return ErrInvalidLogLevel
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

	return nil
}
