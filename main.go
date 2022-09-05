package main

import (
	"ebiten2arkham/gui"
	"ebiten2arkham/renderer"
	"fmt"
	"github.com/HaBaLeS/arkham-go/command"
	"github.com/HaBaLeS/arkham-go/engine"
	"github.com/HaBaLeS/arkham-go/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

var CARDS_JSON = "../data/all_pretty.json"

type Game struct {

	//things that update and Draw
	cardSprites     []*renderer.CardSprite
	guiSprites      []*renderer.GuiSprite
	dialog          *gui.Dialog
	agendaAndAct    *gui.AgendaAndAct
	investigatorGui *renderer.InvestigatorGui

	commandQueue chan command.GuiCommand
	shader       *ebiten.Shader
	fc           float64
}

func (g *Game) Update() error {
	g.fc++
	select {
	case cmd := <-g.commandQueue:
		g.handleCommand(cmd)
	default:
		//do nothing for unblocking the command
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	mx, my := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
	inputConsumed := false

	if g.dialog != nil {
		inputConsumed = true
		g.dialog.Update()
	}

	g.agendaAndAct.Update(float64(mx), float64(my), clicked)
	for _, cs := range g.guiSprites {
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
	for _, cs := range g.cardSprites {
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

func fireEventButtonClicked(cs *renderer.GuiSprite) {
	log.Printf("Button click: %s", cs.Id)
	cs.OnClickFunc()
}

func fireEventCardClicked(cs *renderer.CardSprite) {
	log.Printf("Card click: (%s): %s", cs.Card().CardCode(), cs.Card().Base().Name)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, cs := range g.cardSprites {
		cs.Draw(screen)
	}

	w, h := screen.Size()
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Time":       float32(g.fc) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"ScreenSize": []float32{float32(w), float32(h)},
	}

	g.agendaAndAct.Draw(screen)
	for _, gs := range g.guiSprites {
		gs.Draw(screen)
	}

	if g.dialog != nil {
		screen.DrawRectShader(1920, 1080, g.shader, op)
		g.dialog.Draw(screen)
	}

	mx, my := ebiten.CursorPosition()
	out := fmt.Sprintf("Mouse: %f, %f\n In: %f", float64(mx), float64(my), g.fc)
	ebitenutil.DebugPrint(screen, out)
}

func (g *Game) Layout(ow, oh int) (sh, sw int) {
	return 1920, 1080
}

func main() {

	gui.InitFonts()

	game := &Game{}
	game.init()

	ebiten.SetWindowTitle("Arkham-go")
	ebiten.SetFullscreen(false)
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) init() {

	//Create Game
	g.commandQueue = make(chan command.GuiCommand, 100)

	runtime.Init(g.commandQueue)

	g.cardSprites = make([]*renderer.CardSprite, 0)
	g.guiSprites = make([]*renderer.GuiSprite, 0)

	//fixme load all ressouces
	// scenario, player, deck etc
	for _, p := range runtime.ScenarioSession().Player {
		g.InitCardSpritesForDeck(p)
	}

	//Start Game loop (async -- waits for user input)
	engine := engine.BuildArkhamGame()
	go engine.Start()

	//Render Setup after here
	oneCard := renderer.NewCardSprite(runtime.ScenarioSession().StartLocation)
	oneCard.X = 1920/2 - 300/2
	oneCard.Y = 1080/2 - 419/2
	oneCard.Card().Base().Flipped = true
	oneCard.Enable()

	g.cardSprites = append(g.cardSprites, oneCard)

	//Load investigation Phase GUI
	g.investigatorGui = &renderer.InvestigatorGui{}
	g.guiSprites = append(g.guiSprites, g.investigatorGui.LoadGuiSprites()...)

	g.shader = renderer.GetShader()

	g.agendaAndAct = gui.NewAgendaAndAct()
}

func (g *Game) InitCardSpritesForDeck(deck *runtime.PlayerDeck) {
	sprite := renderer.NewCardSprite(deck.Investigator)
	g.cardSprites = append(g.cardSprites, sprite)

	for _, v := range deck.Cards {
		sprite := renderer.NewCardSprite(v)
		g.cardSprites = append(g.cardSprites, sprite)
	}
}

/*
lol
 -> Gui needs a command thing

 -> Ask XXXXX
 -> Choose YYYYY
 -> Wait for input
 -> Confirm ....

 Big question is how we place cards!! ... i thing we need play-areas where cards are playes in a grid
 We need a Zoom in Card ting

*/

func (g *Game) handleCommand(cmd command.GuiCommand) {
	switch x := cmd.(type) {
	case *command.InfoCommand:
		log.Printf("Reveiced: Info: %v", cmd)
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
		log.Panicf("Unknown GuiCommand %v, %t", cmd, cmd)
	}
}

func (g *Game) playCard(cmd *command.PlayCardCommand) {
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

func (g *Game) getCardSprite(ccode string) *renderer.CardSprite {
	for _, v := range g.cardSprites {
		if v.Card().CardCode() == ccode {
			return v
		}
	}
	return nil
}

func (g *Game) getGuiSprite(id string) *renderer.GuiSprite {
	for _, v := range g.guiSprites {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (g *Game) getGui(s string) *renderer.GuiSprite {
	for _, v := range g.guiSprites {
		if v.Id == s {
			return v
		}
	}
	return nil
}

func (g *Game) enable(what string) {
	switch what {
	case "investigator_gui":
		g.investigatorGui.Enable()
	default:
		log.Panicf("Do not know what to enable: %s", what)
	}
}

func (g *Game) disable(what string) {
	g.getGuiSprite(what).Hidden()
}

func (g *Game) showDialog(x *command.DecisionDialog) {
	//disable all Inputhandling but the dialog

	//show grey out overlay

	g.dialog = gui.NewDialog(x.Question, x.Options)

}

func (g *Game) removeDialog() {
	g.dialog = nil
}

func (g *Game) showCardDialog(x *command.ShowCardDialog) {
	g.dialog = gui.NewCardDialog(x.CardCode, x.ButtonText, x.Front)
}
