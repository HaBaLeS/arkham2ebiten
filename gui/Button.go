package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

type Button struct {
	X, Y, W, H float64
	text       string
	Enabled    bool
}

func (b *Button) Draw(target *ebiten.Image) {
	if b.Enabled {
		ebitenutil.DrawRect(target, b.X, b.Y, b.W, b.H, colornames.Greenyellow)
		text.Draw(target, b.text, normalFnt, int(b.X+5), int(b.Y+22), colornames.Black)
	}
}

func (b *Button) Contains(x, y float64) bool {
	if x > b.X && x < b.X+b.W && y > b.H && y < b.Y+b.H {
		return true
	}
	return false
}
