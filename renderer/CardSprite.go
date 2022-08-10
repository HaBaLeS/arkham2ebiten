package renderer

import (
	"github.com/HaBaLeS/arkham-go/card"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"os"
	"path"
)

type CardSprite struct {
	card         card.ArkhamCard
	enabled      bool
	X, Y         float64
	rotation     float64
	drawable     *ebiten.Image
	drawableBack *ebiten.Image
	width        float64
	height       float64
	Greyout      bool
	Scale        float64
}

func NewCardSprite(card card.ArkhamCard) *CardSprite {
	cs := &CardSprite{
		card: card,
	}
	cs.init()
	return cs
}

func (cs *CardSprite) addFront(file string) {
	fh, err := os.Open(path.Join("../data/leech-img2/", file))
	if err != nil {
		panic(err)
	}
	i, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}
	cs.drawable = ebiten.NewImageFromImage(i)

	w, h := cs.drawable.Size()
	cs.width = float64(w)
	cs.height = float64(h)
}

func (cs *CardSprite) addBack(file string) {
	fh, err := os.Open(path.Join("../data/leech-img2/", file))
	if err != nil {
		panic(err)
	}
	i, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}
	cs.drawableBack = ebiten.NewImageFromImage(i)
}

func (cs *CardSprite) Draw(screen *ebiten.Image) {
	if !cs.enabled {
		return
	}
	op := &ebiten.DrawImageOptions{}

	if cs.card.Base().Tapped {
		cs.rotation = math.Pi / 2
	}

	op.GeoM.Scale(cs.Scale, cs.Scale)
	op.GeoM.Translate(cs.X, cs.Y)

	if cs.Greyout {
		op.ColorM.ChangeHSV(1, 0.5, 1)
	}

	op.Filter = ebiten.FilterLinear

	if cs.card.Base().Flipped {
		screen.DrawImage(cs.drawable, op)
	} else {
		screen.DrawImage(cs.drawableBack, op)
	}
}

func (cs *CardSprite) init() {
	cs.addFront(cs.card.Base().Image)
	cs.addBack(cs.card.Base().BackImage)
	cs.enabled = true
	cs.Greyout = false
	cs.Scale = 1
}

// If you implement small cards (scalled down) or zoom or similar, make sure you update the collision detection too
func (cs *CardSprite) Contains(x, y float64) bool {
	if x > cs.X && x < cs.X+cs.width*cs.Scale && y > cs.Y && y < cs.Y+cs.height*cs.Scale {
		return true
	}
	return false
}

func (cs *CardSprite) Card() card.ArkhamCard {
	return cs.card
}
