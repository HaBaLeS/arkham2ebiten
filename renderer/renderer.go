package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Drawable interface {
	Draw(target *ebiten.Image)
}

type Updatable interface {
	Update()
}

type Toggleable interface {
	Enable()
	Disable()
}

type Hideable interface {
	Show()
	Hide()
}
