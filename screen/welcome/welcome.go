package welcome

import (
	"ebiten2arkham/audio"
	"ebiten2arkham/gui"
	"ebiten2arkham/gui/button"
	"ebiten2arkham/input"
	"ebiten2arkham/renderer"
	"ebiten2arkham/screen"
	"fmt"
	"github.com/HaBaLeS/arkham-go/command"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"os"
)

type WelcomeScreen struct {
	bgBgImage *ebiten.Image
	input     *input.Input
	drawable  []renderer.Drawable
	updatable []renderer.Updatable
	bgSound   *audio.BackgroundMusicPlayer
}

func NewScreen() screen.Screen {

	ws := &WelcomeScreen{
		bgBgImage: renderer.LoadImage("bg/1631463.jpg"),
		drawable:  make([]renderer.Drawable, 0),
		updatable: make([]renderer.Updatable, 0),
		bgSound:   audio.NewBackgroundMusicPlayer("Come-Play-with-Me.mp3"),
	}
	ws.init()
	return ws
}

var frame = 0

func (s *WelcomeScreen) Resume() {
	s.bgSound.PlayLoop()
}

func (s *WelcomeScreen) Pause() {
	s.bgSound.StopLoop()
}

func (s *WelcomeScreen) Update() error {
	for _, u := range s.updatable {
		u.Update()
	}

	//onFrame60Start
	frame++
	if frame == 60 {
		s.bgSound.PlayLoop()
	}
	return nil
}

func (s *WelcomeScreen) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.bgBgImage, nil)
	for _, u := range s.drawable {
		u.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f,%f", input.X, input.Y))
}

func (s *WelcomeScreen) init() {

	logo := gui.NewText(100, 50, 350, "Arkham Horror")
	logo.Shadowed = true
	logo.Enable()
	s.addDrawable(logo)

	startButton := button.NewButton(1200, 400, 125, "start")
	startButton.Enable()
	startButton.DefaultTextColor = color.RGBA{47, 79, 79, 50}
	startButton.ButtonClickFunc = func() {
		command.SendEbitenCommand(&command.SwitchScreen{
			TargetScreen: "game",
		})
	}
	startButton.MouseOverFunc = func() {
		startButton.FrameTextColor = color.RGBA{10, 10, 10, 200}

	}
	s.addUpdatable(startButton)
	s.addDrawable(startButton)

	settingsButton := button.NewButton(1200, 480, 125, "options")
	settingsButton.Enable()
	settingsButton.DefaultTextColor = color.RGBA{47, 79, 79, 50}
	settingsButton.ButtonClickFunc = func() {
		command.SendEbitenCommand(&command.SwitchScreen{
			TargetScreen: "settings",
		})
	}
	settingsButton.MouseOverFunc = func() {
		settingsButton.FrameTextColor = color.RGBA{10, 10, 10, 200}
	}
	s.addUpdatable(settingsButton)
	s.addDrawable(settingsButton)

	exitButton := button.NewButton(1200, 620, 125, "exit")
	exitButton.Enable()
	exitButton.DefaultTextColor = color.RGBA{47, 79, 79, 50}
	exitButton.ButtonClickFunc = func() {
		fmt.Println("Clicked Exit")
		os.Exit(-1)
	}
	exitButton.MouseOverFunc = func() {
		exitButton.FrameTextColor = color.RGBA{10, 10, 10, 200}
	}

	s.addUpdatable(exitButton)
	s.addDrawable(exitButton)

}

func (s *WelcomeScreen) addDrawable(d renderer.Drawable) {
	//fixme ... reserve an amount and keep an index, so thart
	//Zero allocation pollicy!!
	s.drawable = append(s.drawable, d)
}

func (s *WelcomeScreen) addUpdatable(u renderer.Updatable) {
	//fixme ... reserve an amount and keep an index, so thart
	//Zero allocation pollicy!!
	s.updatable = append(s.updatable, u)
}
