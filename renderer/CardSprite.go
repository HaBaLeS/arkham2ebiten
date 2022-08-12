package renderer

import (
	"github.com/HaBaLeS/arkham-go/card"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image"
	"log"
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
	SubImage     image.Rectangle
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
		log.Panicf("Could not load: %s. error:%s", file, err)
	}
	i, _, err := image.Decode(fh)
	if err != nil {
		log.Panicf("Could not load: %s. error:%s", file, err)
	}
	cs.drawable = ebiten.NewImageFromImage(i)

	w, h := cs.drawable.Size()
	cs.width = float64(w)
	cs.height = float64(h)
}

func (cs *CardSprite) addBack(file string) {
	if file == "" {
		return
	}
	fh, err := os.Open(path.Join("../data/leech-img2/", file))
	if err != nil {
		log.Panicf("Could not load: %s. error:%s", file, err)
	}
	i, _, err := image.Decode(fh)
	if err != nil {
		log.Panicf("Could not load: %s. error:%s", file, err)
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

	op.Filter = ebiten.FilterLinear //filter linear always ... this is better for readability

	toDraw := cs.drawable
	if cs.SubImage.Max.X != 0 {
		si := toDraw.SubImage(cs.SubImage)
		toDraw = ebiten.NewImageFromImage(si)
	} else if cs.card.Base().Flipped {
		toDraw = cs.drawableBack
	}
	screen.DrawImage(toDraw, op)

	//Fixme add  to extra render thing !? .. or not ?
	if cs.card.CardType() == card.LocationType {
		loc := card.AcAsLocation(cs.card)
		if loc.ActiveClueTokens > 0 {
			for i := 0; i < loc.ActiveClueTokens; i++ {
				ebitenutil.DrawRect(screen, cs.X+60*float64(i)+30, cs.Y+cs.height-100, 50, 50, colornames.LightGreenA700)
			}
		}
	}
}

func (cs *CardSprite) init() {
	cs.addFront(cs.card.Base().Image)
	cs.addBack(cs.card.Base().BackImage)
	cs.enabled = false
	cs.Greyout = false
	cs.Scale = 1
}

// If you implement small cards (scaled down) or zoom or similar, make sure you update the collision detection too
func (cs *CardSprite) Contains(x, y float64) bool {
	if x > cs.X && x < cs.X+cs.width*cs.Scale && y > cs.Y && y < cs.Y+cs.height*cs.Scale {
		return true
	}
	return false
}

func (cs *CardSprite) Card() card.ArkhamCard {
	return cs.card
}

func (cs *CardSprite) Enable() {
	cs.enabled = true
}

func (cs *CardSprite) Disable() {
	cs.enabled = false
}
