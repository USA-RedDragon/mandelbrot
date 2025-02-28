package game

import (
	"fmt"

	"github.com/USA-RedDragon/mandelbrot/internal/mandelbrot"
	"github.com/USA-RedDragon/mandelbrot/internal/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	mandelbrot *mandelbrot.Mandelbrot
	width      uint
	height     uint
	ui         *ebitenui.UI
	exit       bool
}

func NewGame(width, height uint) (*Game, error) {
	ebiten.SetWindowSize(int(width), int(height))
	ebiten.SetWindowTitle("Fractal Explorer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	eui := &ebitenui.UI{
		Container: root,
	}

	res, err := ui.LoadResources()
	if err != nil {
		return nil, fmt.Errorf("error loading resources: %w", err)
	}

	game := &Game{
		mandelbrot: mandelbrot.NewMandelbrot(int(width), int(height)),
		width:      width,
		height:     height,
		ui:         eui,
		exit:       false,
	}

	manager := NewUIManager(game)
	ui.CreateToolbar(manager, eui, res)

	return game, nil
}

func (g *Game) Update() error {
	if g.exit {
		return ebiten.Termination
	}

	g.ui.Update()
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		x, y := ebiten.CursorPosition()

		desiredCursorPoint := g.mandelbrot.ScreenToViewport(x, y)
		g.mandelbrot.Scale(1 + -wheelY*0.1)
		cursorPointAfterScale := g.mandelbrot.ScreenToViewport(x, y)

		// Recenter based on the difference between the desired point and the cursor point
		g.mandelbrot.Center(g.mandelbrot.GetCenter() + (desiredCursorPoint - cursorPointAfterScale))
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mandelbrot.Update()
	screen.WritePixels(g.mandelbrot.Framebuffer)
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.mandelbrot.Relayout(outsideWidth, outsideHeight)

	return outsideWidth, outsideHeight
}
