package mandelbrot

import (
	"math/cmplx"
	"sync"
)

type Mandelbrot struct {
	width, height int
	framebuffer   []byte
	maxIterations int64
	needsUpdate   bool
	scale         float64
	center        complex128
	exponent      complex128
	startingZ     complex128
	startingC     complex128
	julia         bool
}

type MandelbrotPixel struct {
	X     int
	Y     int
	Color [4]byte
}

const (
	boundMinX = -2
	boundMaxX = 1
	boundMinY = -1
	boundMaxY = 1
)

func NewMandelbrot(width, height int) *Mandelbrot {
	return &Mandelbrot{
		width:         width,
		height:        height,
		framebuffer:   make([]byte, width*height*4),
		maxIterations: 1000,
		needsUpdate:   true,
		scale:         1,
		center:        complex(0, 0),
		exponent:      complex(2, 0),
		startingZ:     complex(0, 0),
		startingC:     complex(-0.63, 0.34),
		julia:         false,
	}
}

func (m *Mandelbrot) GetFramebuffer() []byte {
	return m.framebuffer
}

func (m *Mandelbrot) GetExponent() complex128 {
	return m.exponent
}

func (m *Mandelbrot) GetStartingZ() complex128 {
	return m.startingZ
}

func (m *Mandelbrot) GetStartingC() complex128 {
	return m.startingC
}

func (m *Mandelbrot) Reset() {
	m.scale = 1
	m.center = complex(0, 0)
	m.exponent = complex(2, 0)
	m.startingZ = complex(0, 0)
	m.startingC = complex(-0.63, 0.34)
	m.julia = false
	m.needsUpdate = true
}

func (m *Mandelbrot) SetStartingC(startingC complex128) {
	if m.startingC == startingC {
		return
	}
	m.startingC = startingC
	m.needsUpdate = true
}

func (m *Mandelbrot) SetStartingZ(startingZ complex128) {
	if m.startingZ == startingZ {
		return
	}
	m.startingZ = startingZ
	m.needsUpdate = true
}

func (m *Mandelbrot) SetExponent(exponent complex128) {
	if m.exponent == exponent {
		return
	}
	m.exponent = exponent
	m.needsUpdate = true
}

func (m *Mandelbrot) ScaleBy(factor float64) {
	newscale := m.scale * factor
	if newscale > 1 {
		newscale = 1
	}
	if newscale == m.scale {
		return
	}
	m.scale = newscale
	m.needsUpdate = true
}

func (m *Mandelbrot) Scale(scale float64) {
	if scale == m.scale {
		return
	}
	m.scale = scale
	m.needsUpdate = true
}

func (m *Mandelbrot) GetCenter() complex128 {
	return m.center
}

func (m *Mandelbrot) SetJulia(julia bool) {
	if m.julia == julia {
		return
	}
	m.julia = julia
	m.needsUpdate = true
}

func (m *Mandelbrot) IsJulia() bool {
	return m.julia
}

func (m *Mandelbrot) viewport() [4]float64 {
	return [4]float64{
		boundMinX*m.scale + real(m.center),
		boundMinY*m.scale + imag(m.center),
		boundMaxX*m.scale + real(m.center),
		boundMaxY*m.scale + imag(m.center),
	}
}

func (m *Mandelbrot) ViewportToScreen(point complex128) (x, y int) {
	vp := m.viewport()

	x = int((real(point)-vp[0])/(vp[2]-vp[0])*float64(m.width)) + 1
	y = int((imag(point)-vp[1])/(vp[3]-vp[1])*float64(m.height)) + 1

	return
}

func (m *Mandelbrot) ScreenToViewport(x, y int) complex128 {
	vp := m.viewport()

	real := float64(x)/float64(m.width)*vp[2] + (1-float64(x)/float64(m.width))*vp[0]
	imag := float64(y)/float64(m.height)*vp[3] + (1-float64(y)/float64(m.height))*vp[1]

	return complex(real, imag)
}

func (m *Mandelbrot) Center(center complex128) {
	m.center = center
	m.needsUpdate = true
}

func (m *Mandelbrot) Update() {
	if !m.needsUpdate {
		return
	}
	m.needsUpdate = false

	pixelChan := make(chan MandelbrotPixel, m.width*m.height)
	fbWG := sync.WaitGroup{}
	fbWG.Add(1)
	go func() {
		defer fbWG.Done()
		count := 0
		for pixel := range pixelChan {
			i := (pixel.Y*m.width + pixel.X) * 4
			copy(m.framebuffer[i:i+4], pixel.Color[:])
			count++
			if count == m.width*m.height {
				close(pixelChan)
				return
			}
		}
	}()

	pixelWG := sync.WaitGroup{}
	for y := range m.height {
		for x := range m.width {
			pixelWG.Add(1)
			go func(x, y int) {
				defer pixelWG.Done()
				var z complex128
				var c complex128
				if m.julia {
					z = m.ScreenToViewport(x, y)
					c = m.startingC
				} else {
					z = m.startingZ
					c = m.ScreenToViewport(x, y)
				}

				pixelChan <- m.mandelbrot(x, y, z, m.exponent, c)
			}(x, y)
		}
	}
	pixelWG.Wait()
	fbWG.Wait()
}

func (m *Mandelbrot) mandelbrot(x, y int, z complex128, exponent complex128, c complex128) MandelbrotPixel {
	pixel := [4]byte{0, 0, 0, 255}
	brot := MandelbrotPixel{X: x, Y: y, Color: pixel}
	n := int64(0)

	for n < m.maxIterations && cmplx.Abs(z) < 2 {
		z = cmplx.Pow(z, exponent) + c
		n++
	}

	if n == m.maxIterations {
		return brot
	}

	brot.Color = m.color(n)
	return brot
}

func (m *Mandelbrot) Relayout(width, height int) {
	if m.width == width && m.height == height {
		return
	}
	m.width = width
	m.height = height
	m.framebuffer = make([]byte, width*height*4)
	m.needsUpdate = true
}
