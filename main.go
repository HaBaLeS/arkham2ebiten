package main

import (
	"ebiten2arkham/config"
	"ebiten2arkham/gui/font"
	"ebiten2arkham/input"
	"ebiten2arkham/screen"
	"ebiten2arkham/screen/experiment"
	"ebiten2arkham/screen/game"
	"ebiten2arkham/screen/settings"
	"ebiten2arkham/screen/welcome"
	"github.com/HaBaLeS/arkham-go/command"
	"github.com/HaBaLeS/arkham-go/engine"
	"github.com/HaBaLeS/arkham-go/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

type Game struct {
	commandQueue  chan command.EbitenCommand
	currentScreen screen.Screen
	screenMap     map[string]screen.Screen
	input         *input.Input
}

func (g *Game) Update() error {

	g.currentScreen.Update()

	select {
	case cmd := <-g.commandQueue:
		g.handleCommand(cmd)
	default:
		//do nothing for unblocking the command
	}

	g.input.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.currentScreen.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (sh, sw int) {
	return 1920, 1080
}

func main() {

	//Load config file
	config.Init()

	font.InitFonts()

	game := &Game{}
	game.init()

	ebiten.SetWindowTitle("Arkham-go")
	ebiten.SetFullscreen(config.Cfg.Fullscreen)
	if !config.Cfg.Fullscreen {
		ebiten.SetWindowSize(1920, 1080)
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) init() {

	g.input = &input.Input{}

	//Create Game
	g.commandQueue = make(chan command.EbitenCommand, 100)

	//Start Game loop (async -- waits for user input)
	engine := engine.BuildArkhamGame()
	runtime.Init(g.commandQueue)

	g.screenMap = make(map[string]screen.Screen)
	g.screenMap["main"] = welcome.NewScreen()
	g.screenMap["settings"] = settings.NewScreen()
	g.screenMap["game"] = game.NewScreen()
	g.screenMap["experiment"] = experiment.NewScreen()

	g.currentScreen = g.screenMap["experiment"]

	go engine.Start()

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
	case *command.SwitchScreen:
		g.changeScreen(x.TargetScreen)
	default:
		//Did you use a pointer when sending the command?
		log.Panicf("Unknown GuiCommand %v, %t\n Did you send a pointer?", x, cmd)
	}
}

func (g *Game) changeScreen(targetScreen string) {
	g.currentScreen.Pause()
	g.currentScreen = g.screenMap[targetScreen] //fixme use own type for screens!
	g.currentScreen.Resume()
}
