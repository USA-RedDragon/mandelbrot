package game

// UIManager facilitates communication between the game and the UI
type UIManager struct {
	game *Game
}

func NewUIManager(game *Game) *UIManager {
	return &UIManager{
		game: game,
	}
}

func (m *UIManager) Exit() {
	m.game.exit = true
}

func (m *UIManager) Reset() {
	m.game.mandelbrot.Reset()
}

func (m *UIManager) SetExponentReal(exponent float64) {
	m.game.mandelbrot.SetExponent(complex(exponent, 0))
}

func (m *UIManager) SetExponentImag(exponent float64) {
	m.game.mandelbrot.SetExponent(complex(real(m.game.mandelbrot.GetExponent()), exponent))
}

func (m *UIManager) SetStartingZReal(z float64) {
	m.game.mandelbrot.SetStartingZ(complex(z, imag(m.game.mandelbrot.GetStartingZ())))
}

func (m *UIManager) SetStartingZImag(z float64) {
	m.game.mandelbrot.SetStartingZ(complex(real(m.game.mandelbrot.GetStartingZ()), z))
}

func (m *UIManager) SetStartingCReal(c float64) {
	m.game.mandelbrot.SetStartingC(complex(c, imag(m.game.mandelbrot.GetStartingC())))
}

func (m *UIManager) SetStartingCImag(c float64) {
	m.game.mandelbrot.SetStartingC(complex(real(m.game.mandelbrot.GetStartingC()), c))
}

func (m *UIManager) IsJulia() bool {
	return m.game.mandelbrot.IsJulia()
}

func (m *UIManager) SetJulia(julia bool) {
	m.game.mandelbrot.SetJulia(julia)
}

func (m *UIManager) SetMaxIterations(iterations uint64) {
	m.game.mandelbrot.SetMaxIterations(iterations)
}
