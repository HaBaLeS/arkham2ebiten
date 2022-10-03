package button

import (
	"ebiten2arkham/gui"
	"ebiten2arkham/input"
	"ebiten2arkham/math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	"image/color"
)

type MouseOverFunc func()
type ButtonClickFunc func()

type Button struct {
	math.Rectangle
	buttonText       string
	textSize         float64
	text             *gui.Text
	BgColor          color.Color
	FocusColor       color.Color
	HasBgColor       bool
	HasFocusColor    bool
	enabled          bool
	focus            bool
	MouseOverFunc    MouseOverFunc
	ButtonClickFunc  ButtonClickFunc
	DefaultTextColor color.Color
	FrameTextColor   color.Color
}

func NewButton(x, y float64, size float64, text string) *Button {
	b := &Button{
		buttonText: text,
		textSize:   size,
	}
	b.X = x
	b.Y = y
	b.init()

	return b
}

func (b *Button) Draw(target *ebiten.Image) {
	if b.enabled {
		if b.HasBgColor && !b.focus {
			ebitenutil.DrawRect(target, b.X, b.Y, b.W, b.H, b.BgColor)
		} else if b.focus && b.HasFocusColor {
			ebitenutil.DrawRect(target, b.X, b.Y, b.W, b.H, b.FocusColor)
		}
		b.text.Draw(target)
	}
}

func (b *Button) Update() {
	if !b.enabled {
		return
	}

	//Reset color
	b.FrameTextColor = b.DefaultTextColor

	if b.Contains(input.X, input.Y) {
		b.focus = true
		if b.MouseOverFunc != nil {
			b.MouseOverFunc()
		}
		if input.LMB {
			if b.ButtonClickFunc != nil {
				b.ButtonClickFunc()
			}
		}
	} else {
		b.focus = false
	}

	b.text.Color = b.FrameTextColor
}

func (b *Button) init() {
	b.text = gui.NewText(b.X, b.Y, b.textSize, b.buttonText)
	b.HasBgColor = false
	b.HasFocusColor = false
	b.DefaultTextColor = colornames.Black
	b.BgColor = colornames.Grey
	b.FocusColor = colornames.Green

	b.W = b.text.W
	b.H = b.text.H
}

func (b *Button) Enable() {
	b.enabled = true
	b.text.Enable()
}
func (b *Button) Disable() {
	b.enabled = false
	b.text.Disable()
}
