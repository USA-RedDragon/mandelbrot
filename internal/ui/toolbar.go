package ui

import (
	goimage "image"
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type Toolbar struct {
	container    *widget.Container
	explorerMenu *widget.Button
	quitButton   *widget.Button
}

func CreateToolbar(manager Manager, ui *ebitenui.UI, res *resources) {
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.Black)),

		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			),
		),

		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{StretchHorizontal: true}),
		),
	)

	exponent := newToolbarButton(res, "Exponent")
	var (
		exponentReal = newToolbarNumberEntry(res,
			"Real",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseFloat(newInputText, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if f, err := strconv.ParseFloat(args.InputText, 64); err == nil {
					manager.SetExponentReal(f)
				}
			})
		exponentImag = newToolbarNumberEntry(res,
			"Imag",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseFloat(newInputText, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if f, err := strconv.ParseFloat(args.InputText, 64); err == nil {
					manager.SetExponentImag(f)
				}
			})
	)
	exponent.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, exponentReal, exponentImag)
		}),
	)

	iterations := newToolbarButton(res, "Iterations")
	var (
		iterationsInput = newToolbarNumberEntry(res,
			"Iters",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseUint(newInputText, 10, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if iters, err := strconv.ParseUint(args.TextInput.GetText(), 10, 64); err == nil {
					manager.SetMaxIterations(iters)
				}
			})
	)
	iterations.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, iterationsInput)
		}),
	)

	z := newToolbarButton(res, "Z")
	var (
		zReal = newToolbarNumberEntry(res,
			"Real",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseFloat(newInputText, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if f, err := strconv.ParseFloat(args.InputText, 64); err == nil {
					manager.SetStartingZReal(f)
				}
			})
		zImag = newToolbarNumberEntry(res,
			"Imag",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseFloat(newInputText, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if f, err := strconv.ParseFloat(args.InputText, 64); err == nil {
					manager.SetStartingZImag(f)
				}
			})
	)
	z.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, zReal, zImag)
		}),
	)

	c := newToolbarButton(res, "c")
	var (
		cReal = newToolbarNumberEntry(res,
			"Real",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseFloat(newInputText, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if f, err := strconv.ParseFloat(args.InputText, 64); err == nil {
					manager.SetStartingCReal(f)
				}
			})
		cImag = newToolbarNumberEntry(res,
			"Imag",
			func(newInputText string) (bool, *string) {
				if _, err := strconv.ParseFloat(newInputText, 64); err != nil {
					return false, nil
				}
				return true, &newInputText
			},
			func(args *widget.TextInputChangedEventArgs) {
				if f, err := strconv.ParseFloat(args.InputText, 64); err == nil {
					manager.SetStartingCImag(f)
				}
			})
	)
	c.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, cReal, cImag)
		}),
	)

	explorer := newToolbarButton(res, "Explorer")
	var (
		julia = newToolbarMenuEntryCheckbox(res,
			"Julia",
			func(args *widget.CheckboxChangedEventArgs) {
				if args.State == widget.WidgetChecked {
					manager.SetJulia(true)
					z.GetWidget().Disabled = true
					c.GetWidget().Disabled = false
				} else {
					manager.SetJulia(false)
					z.GetWidget().Disabled = false
					c.GetWidget().Disabled = true
				}
			})
		reset = newToolbarMenuEntry(res, "Reset")
		quit  = newToolbarMenuEntry(res, "Quit")
	)
	if manager.IsJulia() {
		z.GetWidget().Disabled = true
	} else {
		c.GetWidget().Disabled = true
	}
	explorer.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, julia, reset, quit)
		}),
	)
	quit.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			manager.Exit()
		}),
	)
	reset.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			manager.Reset()
		}),
	)

	root.AddChild(explorer)
	root.AddChild(iterations)
	root.AddChild(exponent)
	root.AddChild(z)
	root.AddChild(c)

	toolbar := &Toolbar{
		container:    root,
		explorerMenu: explorer,
		quitButton:   quit,
	}
	ui.Container.AddChild(toolbar.container)
}

func newToolbarMenuEntryCheckbox(res *resources, label string, handler widget.CheckboxChangedHandlerFunc) *widget.LabeledCheckbox {
	uncheckedImage := ebiten.NewImage(15, 15)
	uncheckedImage.Fill(color.White)

	checkedImage := ebiten.NewImage(15, 15)
	checkedImage.Fill(color.NRGBA{255, 255, 0, 255})

	buttonImage, _ := loadButtonImage()

	return widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(15, 15),
				),
				widget.ButtonOpts.Image(buttonImage),
				widget.ButtonOpts.DisableDefaultKeys(),
				widget.ButtonOpts.TextPadding(widget.Insets{
					Top:    4,
					Left:   4,
					Right:  32,
					Bottom: 4,
				}),
			),
			widget.CheckboxOpts.Image(&widget.CheckboxGraphicImage{
				Unchecked: &widget.ButtonImageImage{
					Idle: uncheckedImage,
				},
				Checked: &widget.ButtonImageImage{
					Idle: checkedImage,
				},
			}),
			widget.CheckboxOpts.StateChangedHandler(handler),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text(label, res.font, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.White,
		})),
		widget.LabeledCheckboxOpts.Spacing(6),
	)
}

func newToolbarButton(res *resources, label string) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    4,
			Left:   4,
			Right:  32,
			Bottom: 4,
		}),
	)
}

func newToolbarNumberEntry(res *resources, placeholder string, validator widget.TextInputValidationFunc, handler widget.TextInputChangedHandlerFunc) *widget.TextInput {
	face, _ := loadFont(20)
	return widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),
		widget.TextInputOpts.Face(face),
		widget.TextInputOpts.Placeholder(placeholder),
		widget.TextInputOpts.Validation(validator),
		widget.TextInputOpts.SubmitHandler(handler),
	)
}

func newToolbarMenuEntry(res *resources, label string) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ButtonOpts.TextPadding(widget.Insets{Left: 16, Right: 64}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
	)
}

func openToolbarMenu(opener *widget.Widget, ui *ebitenui.UI, entries ...widget.PreferredSizeLocateableWidget) {
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 125})),

		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 1, Bottom: 1}),
			),
		),

		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(64, 0)),
	)

	for _, entry := range entries {
		c.AddChild(entry)
	}

	w, h := c.PreferredSize()

	window := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),

		widget.WindowOpts.CloseMode(widget.CLICK_OUT),

		widget.WindowOpts.Location(
			goimage.Rect(
				opener.Rect.Min.X,
				opener.Rect.Min.Y+opener.Rect.Max.Y,
				opener.Rect.Min.X+w,
				opener.Rect.Min.Y+opener.Rect.Max.Y+opener.Rect.Min.Y+h,
			),
		),
	)

	ui.AddWindow(window)
}
