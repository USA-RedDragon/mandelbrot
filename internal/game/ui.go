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
	m.game.mandelbrot.Center(complex(0, 0))
	m.game.mandelbrot.Scale(1)
	m.game.mandelbrot.SetExponent(complex(2, 0))
}

func (m *UIManager) SetExponent(exponent float64) {
	m.game.mandelbrot.SetExponent(complex(exponent, 0))
}

func (m *UIManager) SetStartingZReal(z float64) {
	startingZ := m.game.mandelbrot.GetStartingZ()
	m.game.mandelbrot.SetStartingZ(complex(z, imag(startingZ)))
}

func (m *UIManager) SetStartingZImag(z float64) {
	startingZ := m.game.mandelbrot.GetStartingZ()
	m.game.mandelbrot.SetStartingZ(complex(real(startingZ), z))
}
