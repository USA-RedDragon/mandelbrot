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
