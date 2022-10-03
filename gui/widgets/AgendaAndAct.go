package widgets

import (
	"ebiten2arkham/gui/button"
	"ebiten2arkham/gui/font"
	"ebiten2arkham/renderer"
	"fmt"
	"github.com/HaBaLeS/arkham-go/card"
	"github.com/HaBaLeS/arkham-go/command"
	"github.com/HaBaLeS/arkham-go/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

type AgendaAndAct struct {
	currentAgendaCard *card.Agenda
	currentActCard    *card.Act
	progressButton    *button.Button

	currentAgendaSprite *renderer.CardSprite
	currentActSprite    *renderer.CardSprite
	Enabled             bool
}

func NewAgendaAndAct() *AgendaAndAct {
	a := &AgendaAndAct{}

	a.currentAgendaCard = runtime.ScenarioSession().CurrentAgenda
	a.currentActCard = runtime.ScenarioSession().CurrentAct

	scale := 0.4
	a.currentActSprite = renderer.NewCardSprite(a.currentActCard)
	a.currentActSprite.X = 1920 - 419*scale
	a.currentActSprite.Y = 5
	a.currentActSprite.Scale = scale
	a.currentActSprite.Enable()

	a.currentAgendaSprite = renderer.NewCardSprite(a.currentAgendaCard)
	a.currentAgendaSprite.X = 1920 - 419*2*scale
	a.currentAgendaSprite.Y = 5
	a.currentAgendaSprite.Scale = scale
	a.currentAgendaSprite.Enable()

	a.progressButton = button.NewButton(a.currentActSprite.X, 200, 30, "Progress Agenda")
	a.progressButton.DefaultTextColor = colornames.Green
	a.progressButton.HasBgColor = true
	a.progressButton.HasFocusColor = true
	a.progressButton.ButtonClickFunc = func() {
		command.SendEngineCommand(&command.ProgressActCommand{})
	}

	return a
}

func (a *AgendaAndAct) Update(mx, my float64, clicked bool) {

	a.progressButton.Update()

	if a.currentActCard.CanProgress() { //fixme dont call this with 20 fps
		a.progressButton.Enable()
	}

}

func (a *AgendaAndAct) Draw(screen *ebiten.Image) {
	a.currentAgendaSprite.Draw(screen)
	a.currentActSprite.Draw(screen)
	a.progressButton.Draw(screen)

	agendaText := fmt.Sprintf("%s\n %d of %d Doom", a.currentAgendaCard.Name, a.currentAgendaCard.ActiveDoom(), a.currentAgendaCard.Doom)
	text.Draw(screen, agendaText, font.NormalFnt, int(a.currentAgendaSprite.X), int(a.currentAgendaSprite.Y+150), colornames.Whitesmoke)

	actText := fmt.Sprintf("%s\n %d of %d Clues", a.currentActCard.Name, a.currentActCard.ActiveClues(), a.currentActCard.Clues)
	text.Draw(screen, actText, font.NormalFnt, int(a.currentActSprite.X), int(a.currentActSprite.Y+150), colornames.Whitesmoke)

}

//fixme do we want a handle command thing?
