package config_test

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/USA-RedDragon/mandelbrot/internal/config"
	"github.com/spf13/pflag"
)

// load runs the loader in an empty temp dir with args as the command line.
func load(t *testing.T, files map[string]string, args ...string) (*config.Config, error) {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(dir)
	fs := pflag.NewFlagSet("mandelbrot", pflag.ContinueOnError)
	c := config.New(fs)
	if err := fs.Parse(args); err != nil {
		t.Fatal(err)
	}
	return c.Load()
}

func mustLoad(t *testing.T, files map[string]string, args ...string) *config.Config {
	t.Helper()
	cfg, err := load(t, files, args...)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return cfg
}

func TestDefaults(t *testing.T) {
	cfg := mustLoad(t, nil)
	want := config.Config{
		LogLevel:      config.LogLevelInfo,
		Width:         720,
		Height:        480,
		MaxIterations: 1000,
		Scale:         1,
		Center:        0,
		Exponent:      2,
		Z:             0,
		C:             complex(-0.63, 0.34),
		Julia:         false,
		Palette:       config.PaletteRainbow,
	}
	if *cfg != want {
		t.Fatalf("got %+v, want %+v", *cfg, want)
	}
}

func TestEnv(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("WIDTH", "1024")
	t.Setenv("HEIGHT", "768")
	t.Setenv("MAX_ITERATIONS", "50")
	t.Setenv("C", "1i")
	t.Setenv("JULIA", "true")
	t.Setenv("PALETTE", "grayscale")
	cfg := mustLoad(t, nil)
	if cfg.LogLevel != config.LogLevelDebug || cfg.Width != 1024 || cfg.Height != 768 ||
		cfg.MaxIterations != 50 || cfg.C != 1i || !cfg.Julia || cfg.Palette != config.PaletteGrayscale {
		t.Fatalf("env not applied: %+v", *cfg)
	}
}

func TestYAMLFile(t *testing.T) {
	yml := "log-level: warn\nwidth: 800\nscale: 0.5\ncenter: \"-0.5+0.25i\"\nexponent: 3\n"
	cfg := mustLoad(t, map[string]string{"config.yaml": yml})
	if cfg.LogLevel != config.LogLevelWarn || cfg.Width != 800 || cfg.Scale != 0.5 ||
		cfg.Center != complex(-0.5, 0.25) || cfg.Exponent != 3 {
		t.Fatalf("file not applied: %+v", *cfg)
	}
}

func TestYMLFileFromFlag(t *testing.T) {
	cfg := mustLoad(t, map[string]string{"other.yml": "height: 600\n"}, "-c", "other.yml")
	if cfg.Height != 600 {
		t.Fatalf("height = %d, want 600", cfg.Height)
	}
}

func TestExplicitMissingFile(t *testing.T) {
	if _, err := load(t, nil, "--config", "missing.yaml"); err == nil {
		t.Fatal("expected an error for a missing --config file")
	}
}

func TestPrecedence(t *testing.T) {
	t.Setenv("WIDTH", "900")
	t.Setenv("HEIGHT", "700")
	cfg := mustLoad(t, map[string]string{"config.yaml": "width: 800\nheight: 600\nmax-iterations: 20\n"},
		"--width=1000")
	if cfg.Width != 1000 {
		t.Errorf("width = %d, want flag value 1000", cfg.Width)
	}
	if cfg.Height != 700 {
		t.Errorf("height = %d, want env value 700", cfg.Height)
	}
	if cfg.MaxIterations != 20 {
		t.Errorf("max-iterations = %d, want file value 20", cfg.MaxIterations)
	}
}

func TestComplexFlags(t *testing.T) {
	t.Setenv("C", "1i")
	cfg := mustLoad(t, nil, "--c=(-0.8+0.156i)", "--z=0.1-0.2i", "--exponent=2.5")
	if cfg.C != complex(-0.8, 0.156) {
		t.Errorf("c = %v, want (-0.8+0.156i)", cfg.C)
	}
	if cfg.Z != complex(0.1, -0.2) {
		t.Errorf("z = %v, want (0.1-0.2i)", cfg.Z)
	}
	if cfg.Exponent != 2.5 {
		t.Errorf("exponent = %v, want 2.5", cfg.Exponent)
	}
}

func TestBadComplex(t *testing.T) {
	if _, err := load(t, nil, "--center=nope"); err == nil {
		t.Fatal("expected a parse error")
	}
}

func TestValidate(t *testing.T) {
	nan := math.NaN()
	inf := math.Inf(1)
	tests := []struct {
		name   string
		mutate func(*config.Config)
		want   error
	}{
		{"log level", func(c *config.Config) { c.LogLevel = "trace" }, config.ErrInvalidLogLevel},
		{"width", func(c *config.Config) { c.Width = 0 }, config.ErrInvalidWidth},
		{"height", func(c *config.Config) { c.Height = 0 }, config.ErrInvalidHeight},
		{"max iterations", func(c *config.Config) { c.MaxIterations = 0 }, config.ErrInvalidMaxIterations},
		{"zero scale", func(c *config.Config) { c.Scale = 0 }, config.ErrInvalidScale},
		{"negative scale", func(c *config.Config) { c.Scale = -1 }, config.ErrInvalidScale},
		{"NaN scale", func(c *config.Config) { c.Scale = nan }, config.ErrInvalidScale},
		{"Inf scale", func(c *config.Config) { c.Scale = inf }, config.ErrInvalidScale},
		{"NaN center", func(c *config.Config) { c.Center = complex(nan, 0) }, config.ErrInvalidCenter},
		{"Inf exponent", func(c *config.Config) { c.Exponent = complex(0, inf) }, config.ErrInvalidExponent},
		{"Inf z", func(c *config.Config) { c.Z = complex(inf, 0) }, config.ErrInvalidZ},
		{"NaN c", func(c *config.Config) { c.C = complex(0, nan) }, config.ErrInvalidC},
		{"palette", func(c *config.Config) { c.Palette = "plaid" }, config.ErrInvalidPalette},
	}
	base := mustLoad(t, nil)
	if err := base.Validate(); err != nil {
		t.Fatalf("defaults should be valid: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := *base
			tt.mutate(&c)
			if err := c.Validate(); !errors.Is(err, tt.want) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestLoadRunsValidate(t *testing.T) {
	_, err := load(t, nil, "--max-iterations=0")
	if !errors.Is(err, config.ErrInvalidMaxIterations) {
		t.Fatalf("Load() = %v, want %v", err, config.ErrInvalidMaxIterations)
	}
	t.Setenv("SCALE", "NaN")
	if _, err := load(t, nil); !errors.Is(err, config.ErrInvalidScale) {
		t.Fatalf("Load() with SCALE=NaN = %v, want %v", err, config.ErrInvalidScale)
	}
	t.Setenv("SCALE", "1")
	t.Setenv("CENTER", "(NaN+1i)")
	if _, err := load(t, nil); !errors.Is(err, config.ErrInvalidCenter) {
		t.Fatalf("Load() with CENTER=(NaN+1i) = %v, want %v", err, config.ErrInvalidCenter)
	}
}
