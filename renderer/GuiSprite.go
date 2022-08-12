package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"os"
	"path"
)

type GuiSprite struct {
	Id           string
	enabled      bool
	X, Y         float64
	rotation     float64
	drawable     *ebiten.Image
	drawableBack *ebiten.Image
	width        float64
	height       float64
	Greyout      bool
	Scale        float64
	OnClickFunc
}

type OnClickFunc func()

func NewGuiSprite(id, file string) *GuiSprite {
	cs := &GuiSprite{
		Id: id,
	}
	cs.addImage(file)
	cs.init()
	return cs
}

func (cs *GuiSprite) addImage(file string) {
	fh, err := os.Open(path.Join("../data/other/", file))
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

func (cs *GuiSprite) Draw(screen *ebiten.Image) {
	if !cs.enabled {
		return
	}
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(cs.Scale, cs.Scale)
	op.GeoM.Translate(cs.X, cs.Y)

	if cs.Greyout {
		op.ColorM.ChangeHSV(1, 1.5, 1)
	}

	op.Filter = ebiten.FilterLinear

	screen.DrawImage(cs.drawable, op)

}

func (cs *GuiSprite) init() {
	cs.enabled = false
	cs.Greyout = false
	cs.Scale = 1
}

// If you implement small cards (scalled down) or zoom or similar, make sure you update the collision detection too
func (cs *GuiSprite) Contains(x, y float64) bool {
	if !cs.enabled {
		return false
	}
	if x > cs.X && x < cs.X+cs.width*cs.Scale && y > cs.Y && y < cs.Y+cs.height*cs.Scale {
		return true
	}
	return false
}

func (cs *GuiSprite) Disable() {
	cs.enabled = false
}

func (cs *GuiSprite) Enable() {
	cs.enabled = true
}

func LoadGuiSprites() []*GuiSprite {
	retList := make([]*GuiSprite, 0)
	retList = append(retList,
		NewGuiSprite("testButton", "button.png"),
		NewGuiSprite("draw_card", "draw_card.png"),
		NewGuiSprite("engage", "engage.png"),
		NewGuiSprite("escape", "engage.png"),
		NewGuiSprite("fight", "fight.png"),
		NewGuiSprite("investigate", "investigate.png"),
		NewGuiSprite("move", "move.png"),
		NewGuiSprite("play_card", "button.png"),
		NewGuiSprite("resource", "resource.png"))

	return retList
}
