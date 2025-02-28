package ui

// Implementation under internal/game/ui.go
type Manager interface {
	Exit()
	Reset()
	SetExponentReal(exponent float64)
	SetExponentImag(exponent float64)
	SetStartingZReal(z float64)
	SetStartingZImag(z float64)
	SetStartingCReal(z float64)
	SetStartingCImag(z float64)
	IsJulia() bool
	SetJulia(julia bool)
}
