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
}
