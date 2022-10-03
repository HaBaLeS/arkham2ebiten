package gui

import (
	"ebiten2arkham/math"
	"github.com/HaBaLeS/arkhamassets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
)

var cache map[string]*opentype.Font

const dpi = 72

type Text struct {
	math.Rectangle
	text        string
	size        float64
	enabled     bool
	Shadowed    bool
	Color       color.Color
	ShadowColor color.Color
	face        font.Face
}

func NewText(x, y float64, size float64, text string) *Text {
	t := &Text{
		size: size,
		text: text,
	}
	t.X = x
	t.Y = y

	t.init()
	return t
}

func (t *Text) Draw(target *ebiten.Image) {
	if t.enabled {
		if t.Shadowed {
			text.Draw(target, t.text, t.face, int(t.X-1), int(t.Y-1+t.H), t.ShadowColor)
		}
		text.Draw(target, t.text, t.face, int(t.X), int(t.Y+t.H), t.Color)
	}
}

func (t *Text) init() {
	var err error
	tt := getFont("data/font/Teutonic.ttf")

	t.face, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    t.size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Panicf("Could not read font", err)
	}
	t.Color = colornames.Black
	t.ShadowColor = colornames.Whitesmoke
	t.W = float64(text.BoundString(t.face, t.text).Dx()) //fixme ... this needs a dirty flag on update text
	t.H = float64(text.BoundString(t.face, t.text).Dy()) // fixme dont call it twice!

}

func (t *Text) Enable() {
	t.enabled = true
}
func (t *Text) Disable() {
	t.enabled = false
}

//this sets a new text -- make sure if you do, that partents get updated too
/*func (t *Text) Set(s string) {
	t.text = s
	t.W = float64(text.BoundString(t.face, t.text).Dx()) //fixme ... this needs a dirty flag on update text
	t.H = float64(text.BoundString(t.face, t.text).Dy()) // fixme dont call it twice!
}*/

func getFont(path string) *opentype.Font {
	if cache == nil {
		cache = make(map[string]*opentype.Font)
	}
	if cache[path] != nil {
		return cache[path]
	}

	fnt, err := arkhamassets.Data.ReadFile(path)
	if err != nil {
		log.Panicf("Could not read font %v", err)
	}

	tt, err := opentype.Parse(fnt)
	if err != nil {
		log.Panicf("Could not read font %v", err)
	}

	cache[path] = tt
	return tt
}
