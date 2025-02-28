package mandelbrot

import (
	"math"

	"goki.dev/cam/hsl"
)

type PaletteMode int

const (
	PaletteModeSimpleGrayscale PaletteMode = iota
	PaletteModeSimpleRainbow
)

type Palette struct {
	mode PaletteMode
}

func NewPalette(mode PaletteMode) *Palette {
	return &Palette{
		mode: mode,
	}
}

func (p *Palette) colorGrayscale(n uint64, maxIterations uint64) [4]byte {
	factor := math.Sqrt(float64(n) / float64(maxIterations))
	intensity := math.Round(float64(maxIterations) * factor)
	color := uint8(intensity * 255 / float64(maxIterations))
	return [4]byte{color, color, color, 255}
}

func (p *Palette) colorRainbow(n uint64, maxIterations uint64) [4]byte {
	factor := float32(n) / float32(maxIterations)
	hue := factor * 360
	r, g, b, a := hsl.New(hue, factor, 0.5).RGBA()
	return [4]byte{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (p *Palette) Color(n uint64, maxIterations uint64) [4]byte {
	switch p.mode {
	case PaletteModeSimpleGrayscale:
		return p.colorGrayscale(n, maxIterations)
	case PaletteModeSimpleRainbow:
		return p.colorRainbow(n, maxIterations)
	default:
		return [4]byte{0, 0, 0, 255}
	}
}
