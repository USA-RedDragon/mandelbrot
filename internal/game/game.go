package game

import (
	"fmt"

	"github.com/USA-RedDragon/mandelbrot/internal"
	"github.com/USA-RedDragon/mandelbrot/internal/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Mandelbrot *internal.Mandelbrot
	Width      int
	Height     int
	UI         *ebitenui.UI
}

func NewGame(width, height int) (*Game, error) {
	ebiten.SetWindowSize(width, height)
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

	toolbar := ui.NewToolbar(eui, res)
	root.AddChild(toolbar.Container)

	return &Game{
		Mandelbrot: internal.NewMandelbrot(width, height),
		Width:      width,
		Height:     height,
		UI:         eui,
	}, nil
}

func (g *Game) Update() error {
	g.UI.Update()
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		x, y := ebiten.CursorPosition()

		desiredCursorPoint := g.Mandelbrot.ScreenToViewport(x, y)
		g.Mandelbrot.Scale(1 + -wheelY*0.1)
		cursorPointAfterScale := g.Mandelbrot.ScreenToViewport(x, y)

		// Recenter based on the difference between the desired point and the cursor point
		g.Mandelbrot.Center(g.Mandelbrot.GetCenter() + (desiredCursorPoint - cursorPointAfterScale))
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Mandelbrot.Update()
	screen.WritePixels(g.Mandelbrot.Framebuffer)
	g.UI.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nTPS: %.2f\nFPS: %.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Mandelbrot.Relayout(outsideWidth, outsideHeight)

	return outsideWidth, outsideHeight
}
