package ui

// Implementation under internal/game/ui.go
type Manager interface {
	Exit()
	Reset()
	SetExponent(exponent float64)
	SetStartingZReal(z float64)
	SetStartingZImag(z float64)
}
