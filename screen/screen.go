package screen

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
	Draw(screen *ebiten.Image)
	Update() error
	Resume()
	Pause()
}
