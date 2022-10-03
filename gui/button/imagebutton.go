package button

import (
	"ebiten2arkham/math"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type ImageButton struct {
	math.Rectangle
	enabled         bool
	icon            *ebiten.Image
	HasFocusColor   bool
	FocusColor      color.Color
	focus           bool
	MouseOverFunc   MouseOverFunc
	ButtonClickFunc ButtonClickFunc
}

func NewImageButton(x, y, scale float64, file string) *ImageButton {
	b := &ImageButton{}

	b.X = x
	b.Y = y

	return b
}

func (b *ImageButton) Draw(target *ebiten.Image) {

}

func (b *ImageButton) Update() {

}

func (b *ImageButton) Enable() {
	b.enabled = true
}
func (b *ImageButton) Disable() {
	b.enabled = false
}
