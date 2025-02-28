package mandelbrot

func (m *Mandelbrot) color(n uint64) [4]byte {
	color := uint8(255 - (n * 255 / m.maxIterations))
	return [4]byte{color, color, color, 255}
}
