package ui

// Implementation under internal/game/ui.go
type Manager interface {
	Exit()
	Reset()
	SetExponent(exponent float64)
}
