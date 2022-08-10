package main

import (
	"ebiten2arkham/renderer"
	"fmt"
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
	cardSprites []*renderer.CardSprite
	guiSprites  []*renderer.GuiSprite
}

func (g *Game) Update() error {
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
	log.Printf("Button Cliecked: %s", cs.Id)
	cs.OnClickFunc()
	cs.Disable()
}

func fireEventCardClicked(cs *renderer.CardSprite) {
	log.Printf("Card Clicked: (%s): %s", cs.Card().CardCode(), cs.Card().Base().Name)
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

	game := initGame()

	ebiten.SetWindowTitle("Arkham-go")
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func initGame() *Game {

	//Create Game
	game := &Game{}
	game.cardSprites = make([]*renderer.CardSprite, 0)
	game.guiSprites = make([]*renderer.GuiSprite, 0)

	//Load CardDB
	db := runtime.NewCardDB()
	err := db.Init(CARDS_JSON)
	if err != nil {
		panic(err)
	}

	//Create Scenario
	scnData := runtime.GetFirstScenarioData(db)
	game.engine = arkham_game.BuildArkhamGame(scnData)

	//Start Game loop (async -- waits for user input)
	go game.engine.Start()

	//Render Setup after here
	oneCard := renderer.NewCardSprite(scnData.StartLocation)
	oneCard.X = 1920/2 - 300/2
	oneCard.Y = 1080/2 - 419/2

	act := renderer.NewCardSprite(scnData.CurrentAct)
	agenda := renderer.NewCardSprite(scnData.CurrentAgenda)

	scale := 0.4

	act.X = 1920 - 419*scale
	act.Y = 5
	act.Card().Base().Flipped = true
	act.Scale = scale

	agenda.X = 1920 - 419*2*scale
	agenda.Y = 5
	agenda.Card().Base().Flipped = true
	agenda.Scale = scale

	game.cardSprites = append(game.cardSprites, oneCard, agenda, act)

	btn := renderer.NewGuiSprite("testButton", "button.png")
	btn.X = 500
	btn.Y = 500 + 425
	btn.OnClickFunc = game.engine.GameStart.Callback //magic trick, bring the callback function form extern
	game.guiSprites = append(game.guiSprites, btn)

	return game
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
