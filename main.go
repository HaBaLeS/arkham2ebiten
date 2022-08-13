package main

import (
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
	//scn         *runtime.Scenario
	cardSprites     []*renderer.CardSprite
	guiSprites      []*renderer.GuiSprite
	commandQueue    chan command.GuiCommand
	investigatorGui *renderer.InvestigatorGui
}

func (g *Game) Update() error {

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
	inputConsumed := false
	for _, cs := range g.guiSprites {
		if cs.Contains(float64(mx), float64(my)) && !inputConsumed {
			cs.Greyout = true
			inputConsumed = true
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
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
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
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
	cs.Disable() //fixme this was only for the first button, no the investigator GUi is disapearing
}

func fireEventCardClicked(cs *renderer.CardSprite) {
	log.Printf("Card click: (%s): %s", cs.Card().CardCode(), cs.Card().Base().Name)
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, cs := range g.cardSprites {
		cs.Draw(screen)
	}

	for _, gs := range g.guiSprites {
		gs.Draw(screen)
	}

	mx, my := ebiten.CursorPosition()
	out := fmt.Sprintf("Mouse: %f, %f", float64(mx), float64(my))
	ebitenutil.DebugPrint(screen, out)
}

func (g *Game) Layout(ow, oh int) (sh, sw int) {
	return 1920, 1080
}

func main() {

	game := &Game{}
	game.init()

	ebiten.SetWindowTitle("Arkham-go")
	ebiten.SetFullscreen(true)

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

	scale := 0.4

	act := renderer.NewCardSprite(runtime.ScenarioSession().CurrentAct)
	act.X = 1920 - 419*scale
	act.Y = 5
	act.Scale = scale
	act.Enable()

	agenda := renderer.NewCardSprite(runtime.ScenarioSession().CurrentAgenda)
	agenda.X = 1920 - 419*2*scale
	agenda.Y = 5
	agenda.Scale = scale
	agenda.Enable()

	g.cardSprites = append(g.cardSprites, oneCard, agenda, act)

	//btn.OnClickFunc = engine.GameStart.Callback //magic trick, bring the callback function form extern
	btn := renderer.NewGuiSprite("testButton", "button.png")
	g.guiSprites = append(g.guiSprites, btn)

	startButton := g.getGui("testButton")
	startButton.OnClickFunc = engine.GameStart.Callback
	startButton.Enable()
	startButton.X = 1920/2 - 100
	startButton.Y = 1080/2 + 419/2 + 5

	//Load investigation Phase GUI
	g.investigatorGui = &renderer.InvestigatorGui{}
	g.guiSprites = append(g.guiSprites, g.investigatorGui.LoadGuiSprites()...)

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
	default:
		log.Panicf("Unknown GioCommand %v", cmd)
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
