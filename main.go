package main

import (
	"ebiten2arkham/renderer"
	"fmt"
	"github.com/HaBaLeS/arkham-go/command"
	arkham_game "github.com/HaBaLeS/arkham-go/modules/arkham-game"
	"github.com/HaBaLeS/arkham-go/modules/gpbge"
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
	engine *gpbge.PhaseEngine
	//scn         *runtime.Scenario
	cardSprites  []*renderer.CardSprite
	guiSprites   []*renderer.GuiSprite
	commandQueue chan command.GuiCommand
}

func (g *Game) Update() error {

	select {
	case cmd := <-g.commandQueue:
		log.Printf("Received Command: %v", cmd)
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
	cs.Disable()
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

	g.cardSprites = make([]*renderer.CardSprite, 0)
	g.guiSprites = make([]*renderer.GuiSprite, 0)

	//Load CardDB
	db := runtime.NewCardDB()
	err := db.Init(CARDS_JSON)
	if err != nil {
		panic(err)
	}

	//Create Scenario
	scnData := runtime.GetFirstScenarioData(db)
	g.engine = arkham_game.BuildArkhamGame(scnData, g.commandQueue)

	//Add players
	d1, err := runtime.LoadPlayerDeckFromFile("../data/decks/deck1.txt", db)
	if err != nil {
		panic(err)
	}
	d2, err := runtime.LoadPlayerDeckFromFile("../data/decks/deck2.txt", db)
	if err != nil {
		panic(err)
	}
	g.engine.AddPlayer(d1)
	g.engine.AddPlayer(d2)

	//fixme load all ressouces
	// scenario, player, deck etc
	g.InitCardSpritesForDeck(d1)
	g.InitCardSpritesForDeck(d2)

	//Start Game loop (async -- waits for user input)
	go g.engine.Start()

	//Render Setup after here
	oneCard := renderer.NewCardSprite(scnData.StartLocation)
	oneCard.X = 1920/2 - 300/2
	oneCard.Y = 1080/2 - 419/2
	oneCard.Card().Base().Flipped = true
	oneCard.Enable()

	scale := 0.4

	act := renderer.NewCardSprite(scnData.CurrentAct)
	act.X = 1920 - 419*scale
	act.Y = 5
	act.Scale = scale
	act.Enable()

	agenda := renderer.NewCardSprite(scnData.CurrentAgenda)
	agenda.X = 1920 - 419*2*scale
	agenda.Y = 5
	agenda.Scale = scale
	agenda.Enable()

	g.cardSprites = append(g.cardSprites, oneCard, agenda, act)

	btn := renderer.NewGuiSprite("testButton", "button.png")
	btn.X = 500
	btn.Y = 500 + 425
	btn.OnClickFunc = g.engine.GameStart.Callback //magic trick, bring the callback function form extern
	g.guiSprites = append(g.guiSprites, btn)

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
