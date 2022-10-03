package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var X float64
var Y float64
var LMB bool

type Input struct {
}

func (i *Input) Update() {
	x, y := ebiten.CursorPosition()
	X = float64(x)
	Y = float64(y)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		LMB = true
	} else {
		LMB = false
	}
}
