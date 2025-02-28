package mandelbrot

import (
	"math"
)

var ()

func (m *Mandelbrot) color(n uint64) [4]byte {
	factor := math.Sqrt(float64(n) / float64(m.maxIterations))
	intensity := math.Round(float64(m.maxIterations) * factor)
	color := uint8(intensity * 255 / float64(m.maxIterations))
	return [4]byte{color, color, color, 255}
}
