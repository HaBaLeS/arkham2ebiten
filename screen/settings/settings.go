package settings

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
)

type SettingsScreen struct {
	bgBgImage *ebiten.Image
	input     *input.Input
	drawable  []renderer.Drawable
	updatable []renderer.Updatable
	bgSound   *audio.BackgroundMusicPlayer
}

func NewScreen() screen.Screen {

	ws := &SettingsScreen{
		bgBgImage: renderer.LoadImage("bg/the_streets_of_innsmouth_by_nemo2d_dehn96w-fullview.jpg"),
		drawable:  make([]renderer.Drawable, 0),
		updatable: make([]renderer.Updatable, 0),
		bgSound:   audio.NewBackgroundMusicPlayer("scott-buckley-i-walk-with-ghosts.mp3"),
	}
	ws.init()
	return ws
}

var frame = 0

func (s *SettingsScreen) Resume() {
	s.bgSound.PlayLoop()
}

func (s *SettingsScreen) Pause() {
	s.bgSound.StopLoop()
}

func (s *SettingsScreen) Update() error {
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

func (s *SettingsScreen) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.bgBgImage, op)
	for _, u := range s.drawable {
		u.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f,%f", input.X, input.Y))

}

func (s *SettingsScreen) init() {

	logo := gui.NewText(100, 50, 200, "Settings")
	logo.Enable()
	s.addDrawable(logo)

	fullscreen := button.NewButton(1200, 400, 125, "Fullscreen")
	fullscreen.Enable()
	fullscreen.DefaultTextColor = color.RGBA{47, 79, 79, 50}
	fullscreen.ButtonClickFunc = func() {
		fmt.Println("ClickedStart")
	}
	fullscreen.MouseOverFunc = func() {
		fullscreen.FrameTextColor = color.RGBA{10, 10, 10, 200}
	}
	s.addUpdatable(fullscreen)
	s.addDrawable(fullscreen)

	mute := button.NewButton(1200, 600, 125, "Mute")
	mute.Enable()
	mute.DefaultTextColor = color.RGBA{47, 79, 79, 50}
	mute.ButtonClickFunc = func() {
		fmt.Println("ClickedStart")
	}
	mute.MouseOverFunc = func() {
		mute.FrameTextColor = color.RGBA{10, 10, 10, 200}
	}
	s.addUpdatable(mute)
	s.addDrawable(mute)

	backButton := button.NewButton(75, 920, 125, "< back")
	backButton.Enable()
	backButton.DefaultTextColor = color.RGBA{47, 79, 79, 50}
	backButton.ButtonClickFunc = func() {
		command.SendEbitenCommand(&command.SwitchScreen{
			"main",
		})
	}
	backButton.MouseOverFunc = func() {
		backButton.FrameTextColor = color.RGBA{10, 10, 10, 200}
	}
	s.addUpdatable(backButton)
	s.addDrawable(backButton)

}

func (s *SettingsScreen) addDrawable(d renderer.Drawable) {
	//fixme ... reserve an amount and keep an index, so thart
	//Zero allocation pollicy!!
	s.drawable = append(s.drawable, d)
}

func (s *SettingsScreen) addUpdatable(u renderer.Updatable) {
	//fixme ... reserve an amount and keep an index, so thart
	//Zero allocation pollicy!!
	s.updatable = append(s.updatable, u)
}
