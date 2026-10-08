// Package config loads the mandelbrot settings from defaults, a config file,
// environment variables and command-line flags.
package config

//go:generate go tool configulator -type Config
//go:generate go tool configulator -type Config -markdown -markdown-file ../../README.md -env-prefix ""
//go:generate go tool configulator -type Config -sample -sample-file ../../config.example.yaml

import (
	"errors"
	"math"
	"math/cmplx"

	configulator "github.com/USA-RedDragon/configulator/v2"
	cpflag "github.com/USA-RedDragon/configulator/v2/flags/pflag"
	"github.com/goccy/go-yaml"
	"github.com/spf13/pflag"
)

// Config is the mandelbrot configuration.
type Config struct {
	LogLevel      LogLevel   `name:"log-level" default:"info" description:"log level: debug, info, warn or error"`
	Width         uint       `name:"width" default:"720" description:"initial window width"`
	Height        uint       `name:"height" default:"480" description:"initial window height"`
	MaxIterations uint64     `name:"max-iterations" default:"1000" description:"maximum iterations per point"`
	Scale         float64    `name:"scale" default:"1" description:"initial zoom scale, smaller zooms in"`
	Center        complex128 `name:"center" default:"0" description:"initial view center"`
	Exponent      complex128 `name:"exponent" default:"2" description:"exponent in z = z^exponent + c"`
	Z             complex128 `name:"z" default:"0" description:"starting z for the Mandelbrot set"`
	C             complex128 `name:"c" default:"(-0.63+0.34i)" description:"c for the Julia set"`
	Julia         bool       `name:"julia" default:"false" description:"start in Julia set mode"`
	Palette       Palette    `name:"palette" default:"rainbow" description:"color palette: rainbow or grayscale"`
}

// LogLevel is the minimum level of log messages to show.
type LogLevel string

// Log levels.
const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Palette is the color palette used to draw points outside the set.
type Palette string

// Palettes.
const (
	PaletteRainbow   Palette = "rainbow"
	PaletteGrayscale Palette = "grayscale"
)

// DefaultConfigPath is the config file read when --config is not given. It
// may be missing.
const DefaultConfigPath = "config.yaml"

// Validation errors.
var (
	ErrInvalidLogLevel      = errors.New("invalid log level")
	ErrInvalidWidth         = errors.New("invalid width")
	ErrInvalidHeight        = errors.New("invalid height")
	ErrInvalidMaxIterations = errors.New("max-iterations must be greater than 0")
	ErrInvalidScale         = errors.New("scale must be a finite number greater than 0")
	ErrInvalidCenter        = errors.New("center must be finite")
	ErrInvalidExponent      = errors.New("exponent must be finite")
	ErrInvalidZ             = errors.New("z must be finite")
	ErrInvalidC             = errors.New("c must be finite")
	ErrInvalidPalette       = errors.New("invalid palette")
)

// Validate checks the loaded config. Load calls it after every layer.
func (c Config) Validate() error {
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

	if c.MaxIterations == 0 {
		return ErrInvalidMaxIterations
	}

	if math.IsNaN(c.Scale) || math.IsInf(c.Scale, 0) || c.Scale <= 0 {
		return ErrInvalidScale
	}

	for _, v := range []struct {
		value complex128
		err   error
	}{
		{c.Center, ErrInvalidCenter},
		{c.Exponent, ErrInvalidExponent},
		{c.Z, ErrInvalidZ},
		{c.C, ErrInvalidC},
	} {
		if cmplx.IsNaN(v.value) || cmplx.IsInf(v.value) {
			return v.err
		}
	}

	switch c.Palette {
	case PaletteRainbow, PaletteGrayscale:
	default:
		return ErrInvalidPalette
	}

	return nil
}

// New builds a loader that reads env vars with no prefix (LOG_LEVEL), the
// YAML file named by -c/--config (config.yaml by default) and the flags in
// fs, which it registers right away.
func New(fs *pflag.FlagSet) *configulator.Configulator[Config] {
	c := configulator.New(ConfigSchema()).
		WithEnvironmentVariables(&configulator.EnvironmentVariableOptions{Prefix: "", Separator: "_"}).
		WithFile(&configulator.FileOptions{
			Search: []string{DefaultConfigPath},
			Decoders: configulator.Decoders{
				".yaml": yaml.Unmarshal,
				".yml":  yaml.Unmarshal,
			},
		})
	return cpflag.Bind(c, fs, ConfigPFlagHooks(), nil)
}
