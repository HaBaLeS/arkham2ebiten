package game

import (
	"ebiten2arkham/audio"
	"ebiten2arkham/gui"
	"ebiten2arkham/gui/widgets"
	"ebiten2arkham/input"
	"ebiten2arkham/renderer"
	"ebiten2arkham/screen"
	"fmt"
	"github.com/HaBaLeS/arkham-go/command"
	"github.com/HaBaLeS/arkham-go/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"os"
)

type GameScreen struct {
	commandQueue chan command.GuiCommand

	bgBgImage *ebiten.Image
	input     *input.Input
	drawable  []renderer.Drawable
	updatable []renderer.Updatable
	bgSound   *audio.BackgroundMusicPlayer

	//things that update and Draw
	cardSprites     []*renderer.CardSprite
	guiSprites      []*renderer.GuiSprite
	agendaAndAct    *widgets.AgendaAndAct
	investigatorGui *renderer.InvestigatorGui
	shader          *ebiten.Shader
	dialog          *gui.Dialog
	fc              float64
}

func NewScreen() screen.Screen {

	ws := &GameScreen{
		bgBgImage:    renderer.LoadImage("bg/1631463.jpg"),
		drawable:     make([]renderer.Drawable, 0),
		updatable:    make([]renderer.Updatable, 0),
		bgSound:      audio.NewBackgroundMusicPlayer("Come-Play-with-Me.mp3"),
		commandQueue: make(chan command.GuiCommand, 100),
	}
	ws.cardSprites = make([]*renderer.CardSprite, 0)
	ws.guiSprites = make([]*renderer.GuiSprite, 0)
	command.SetGuiChannel(ws.commandQueue)
	ws.init()
	return ws
}

func (s *GameScreen) Resume() {
}

func (s *GameScreen) Pause() {
}

func (s *GameScreen) Update() error {
	s.fc++

	select {
	case cmd := <-s.commandQueue:
		s.handleCommand(cmd)
	default:
		//do nothing for unblocking the command
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	mx, my := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
	inputConsumed := false

	if s.dialog != nil {
		inputConsumed = true
		s.dialog.Update()
	}

	s.agendaAndAct.Update(float64(mx), float64(my), clicked)
	for _, cs := range s.guiSprites {
		if cs.Contains(float64(mx), float64(my)) && !inputConsumed {
			cs.Greyout = true
			inputConsumed = true
			if clicked {
				fireEventButtonClicked(cs)
			}
		} else {
			cs.Greyout = false
		}
	}

	//this is the card layer, if you want to overlay, make sure to handle the input first and prevent this block
	for _, cs := range s.cardSprites {
		if cs.Contains(float64(mx), float64(my)) && !inputConsumed {
			cs.Greyout = true
			inputConsumed = true
			if clicked {
				fireEventCardClicked(cs)
			}
		} else {
			cs.Greyout = false
		}
	}
	return nil
}
func (s *GameScreen) Draw(screen *ebiten.Image) {
	for _, cs := range s.cardSprites {
		cs.Draw(screen)
	}

	w, h := screen.Size()
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Time":       float32(s.fc) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"ScreenSize": []float32{float32(w), float32(h)},
	}

	s.agendaAndAct.Draw(screen)
	for _, gs := range s.guiSprites {
		gs.Draw(screen)
	}

	if s.dialog != nil {
		screen.DrawRectShader(1920, 1080, s.shader, op)
		s.dialog.Draw(screen)
	}

	mx, my := ebiten.CursorPosition()
	out := fmt.Sprintf("Mouse: %f, %f\n In: %f", float64(mx), float64(my), s.fc)
	ebitenutil.DebugPrint(screen, out)
}

func (s *GameScreen) init() {

	//Render Setup after here
	oneCard := renderer.NewCardSprite(runtime.ScenarioSession().StartLocation)
	oneCard.X = 1920/2 - 300/2
	oneCard.Y = 1080/2 - 419/2
	oneCard.Card().Base().Flipped = true
	oneCard.Enable()

	s.cardSprites = append(s.cardSprites, oneCard)

	//fixme load all ressouces
	// scenario, player, deck etc
	for _, p := range runtime.ScenarioSession().Player {
		s.InitCardSpritesForDeck(p)
	}

	//Load investigation Phase GUI
	s.investigatorGui = &renderer.InvestigatorGui{}
	s.guiSprites = append(s.guiSprites, s.investigatorGui.LoadGuiSprites()...)

	s.shader = renderer.GetShader()

	s.agendaAndAct = widgets.NewAgendaAndAct()
}

func fireEventButtonClicked(cs *renderer.GuiSprite) {
	log.Printf("Button click: %s", cs.Id)
	cs.OnClickFunc()
}

func fireEventCardClicked(cs *renderer.CardSprite) {
	log.Printf("Card click: (%s): %s", cs.Card().CardCode(), cs.Card().Base().Name)
}

func (s *GameScreen) InitCardSpritesForDeck(deck *runtime.PlayerDeck) {
	sprite := renderer.NewCardSprite(deck.Investigator)
	s.cardSprites = append(s.cardSprites, sprite)

	for _, v := range deck.Cards {
		sprite := renderer.NewCardSprite(v)
		s.cardSprites = append(s.cardSprites, sprite)
	}
}

func (g *GameScreen) handleCommand(cmd command.GuiCommand) {
	switch x := cmd.(type) {
	case *command.PlayCardCommand:
		log.Printf("Reveiced: PlayCardCommand: %v", cmd)
		g.playCard(x)
	case *command.EnableCommand:
		g.enable(x.What)
	case *command.DisableCommand:
		g.disable(x.What)
	case *command.DecisionDialog:
		g.showDialog(x)
	case *command.RemoveDialog:
		g.removeDialog()
	case *command.ShowCardDialog:
		g.showCardDialog(x)
	default:
		//Did you use a pointer when sending the command?
		log.Panicf("Unknown GuiCommand %v, %t\n Did you send a pointer?", x, cmd)
	}
}

func (g *GameScreen) playCard(cmd *command.PlayCardCommand) {
	s := g.getCardSprite(cmd.CardToPlay)
	s.Enable()
	if cmd.Scale != 0 {
		s.Scale = cmd.Scale
	}
	if cmd.X != 0 {
		s.X = cmd.X
	}
	if cmd.Y != 0 {
		s.Y = cmd.Y
	}
	if cmd.SubImage.Max.X != 0 {
		s.SubImage = cmd.SubImage
	}
}

func (g *GameScreen) getCardSprite(ccode string) *renderer.CardSprite {
	for _, v := range g.cardSprites {
		if v.Card().CardCode() == ccode {
			return v
		}
	}
	return nil
}

func (g *GameScreen) getGuiSprite(id string) *renderer.GuiSprite {
	for _, v := range g.guiSprites {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (g *GameScreen) getGui(s string) *renderer.GuiSprite {
	for _, v := range g.guiSprites {
		if v.Id == s {
			return v
		}
	}
	return nil
}

func (g *GameScreen) enable(what string) {
	switch what {
	case "investigator_gui":
		g.investigatorGui.Enable()
	default:
		log.Panicf("Do not know what to enable: %s", what)
	}
}

func (g *GameScreen) disable(what string) {
	g.getGuiSprite(what).Hidden()
}

func (g *GameScreen) showDialog(x *command.DecisionDialog) {
	//disable all Inputhandling but the dialog

	//show grey out overlay

	g.dialog = gui.NewDialog(x.Question, x.Options)

}

func (g *GameScreen) removeDialog() {
	g.dialog = nil
}

func (g *GameScreen) showCardDialog(x *command.ShowCardDialog) {
	g.dialog = gui.NewCardDialog(x.CardCode, x.ButtonText, x.Front)
}
