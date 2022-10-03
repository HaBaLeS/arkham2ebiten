package gui

import (
	"ebiten2arkham/gui/font"
	"ebiten2arkham/renderer"
	"fmt"
	"github.com/HaBaLeS/arkham-go/command"
	"github.com/HaBaLeS/arkham-go/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

type Dialog struct {
	background *ebiten.Image
	question   string
	options    []string
	w, h       float64
	x, y       float64
	enabled    bool
	buttons    []*DialogButton
	card       *renderer.CardSprite
}

func NewDialog(question string, options []string) *Dialog {
	dlg := &Dialog{
		question:   question,
		options:    options,
		x:          1920 / 2,
		y:          50,
		background: renderer.LoadImage("other/dialog.jpg"),
	}
	dlg.enabled = true
	w, h := dlg.background.Size()
	dlg.w = float64(w)
	dlg.h = float64(h)

	dlg.x -= dlg.w / 2

	dlg.buttons = make([]*DialogButton, len(options))
	for i := 0; i < len(options); i++ {
		dlg.buttons[i] = &DialogButton{
			x:        dlg.x + 200,
			y:        dlg.y + 200 + 100*float64(i),
			w:        200,
			h:        75,
			text:     options[i],
			callback: options[i],
		}
	}
	return dlg
}

func NewCardDialog(cardCode, buttonText string, front bool) *Dialog {
	ak := runtime.CardDBG().GetCard(cardCode)
	var cardText string
	if front {
		cardText = fmt.Sprintf("%s\n%s", ak.Base().Name, ak.Base().Flavor)
	} else {
		cardText = fmt.Sprintf("%s\n%s", ak.Base().BackName, ak.Base().BackFlavor)
	}

	dlg := NewDialog(cardText, []string{buttonText})

	//Fixme you should query the renderer for the card!! do not create a new sprite at runtime!
	dlg.card = renderer.NewCardSprite(ak)
	dlg.card.Enable()
	dlg.card.X = 1920/2 - 419/2
	dlg.card.Y = 1080 / 2
	dlg.card.Scale = 1
	dlg.card.Card().Base().Flipped = !front

	return dlg
}

type DialogButton struct { //fixme replace with gui.Button
	text     string
	callback string

	x, y, w, h float64
}

func (b *DialogButton) contains(x, y float64) bool {
	if x > b.x && x < b.x+b.w && y > b.y && y < b.y+b.h {
		return true
	}
	return false
}

func (dlg *Dialog) Draw(target *ebiten.Image) {
	if !dlg.enabled {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(dlg.x, dlg.y)
	target.DrawImage(dlg.background, op)

	msg := fmt.Sprintf("%s", dlg.question)
	text.Draw(target, msg, font.NormalFnt, int(dlg.x+170), int(dlg.y+100), colornames.Whitesmoke)

	for _, b := range dlg.buttons {
		ebitenutil.DrawRect(target, b.x, b.y, b.w, b.h, colornames.Orange)
		text.Draw(target, b.text, font.NormalFnt, int(b.x+10), int(b.y+50), colornames.Black)
	}

	if dlg.card != nil {
		dlg.card.Draw(target)
	}
}

func (dlg *Dialog) Update() {
	mx, my := ebiten.CursorPosition()
	for _, b := range dlg.buttons {
		if b.contains(float64(mx), float64(my)) {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				command.SendEngineCommand(&command.Decision{
					Answer: b.callback,
				})
				command.SendGuiCommand(&command.RemoveDialog{})
			}
		}
	}
}
