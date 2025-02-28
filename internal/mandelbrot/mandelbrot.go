package mandelbrot

import (
	"math/cmplx"
	"sync"
)

type Mandelbrot struct {
	width, height int
	Framebuffer   []byte
	maxIterations int64
	needsUpdate   bool
	scale         float64
	center        complex128
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
		Framebuffer:   make([]byte, width*height*4),
		maxIterations: 1000,
		needsUpdate:   true,
		scale:         1,
		center:        complex(0, 0),
	}
}

func (m *Mandelbrot) Scale(factor float64) {
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

func (m *Mandelbrot) GetCenter() complex128 {
	return m.center
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
		for {
			select {
			case pixel := <-pixelChan:
				i := (pixel.Y*m.width + pixel.X) * 4
				copy(m.Framebuffer[i:i+4], pixel.Color[:])
				count++
				if count == m.width*m.height {
					close(pixelChan)
					return
				}
			}
		}
	}()

	pixelWG := sync.WaitGroup{}
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			pixelWG.Add(1)
			go func(x, y int) {
				defer pixelWG.Done()
				c := m.ScreenToViewport(x, y)
				pixelChan <- m.mandelbrot(x, y, 0, 2, c)
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

func (m *Mandelbrot) color(n int64) [4]byte {
	color := uint8(255 - (n * 255 / m.maxIterations))
	return [4]byte{color, color, color, 255}
}

func (m *Mandelbrot) Relayout(width, height int) {
	if m.width == width && m.height == height {
		return
	}
	m.width = width
	m.height = height
	m.Framebuffer = make([]byte, width*height*4)
	m.needsUpdate = true
}
