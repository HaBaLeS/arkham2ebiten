package audio

import (
	"ebiten2arkham/config"
	"github.com/HaBaLeS/arkhamassets"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"path"
)

type BackgroundMusicPlayer struct {
	ctx        *audio.Context
	loopPlayer *audio.Player
	volume     float64
	mute       bool
}

func NewBackgroundMusicPlayer(file string) *BackgroundMusicPlayer {

	ctx := audio.CurrentContext()
	if ctx == nil {
		ctx = audio.NewContext(48000)
	}

	p := &BackgroundMusicPlayer{
		ctx: ctx,
	}

	src, err := arkhamassets.Data.Open(path.Join("data/sound", file))
	if err != nil {
		panic(err)
	}
	stream, err := mp3.DecodeWithoutResampling(src)
	if err != nil {
		panic(err)
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())

	p.loopPlayer, err = p.ctx.NewPlayer(loop)
	if err != nil {
		panic(err)
	}

	p.volume = p.loopPlayer.Volume()

	p.Mute(config.Cfg.Mute)

	return p
}

func (p *BackgroundMusicPlayer) Mute(m bool) {
	if m {
		p.loopPlayer.SetVolume(0)
	} else {
		p.loopPlayer.SetVolume(p.volume)
	}
}

func (p *BackgroundMusicPlayer) PlayLoop() {
	p.loopPlayer.Play()
}

func (p *BackgroundMusicPlayer) StopLoop() {
	p.loopPlayer.Pause()
}

func (p *BackgroundMusicPlayer) Shutdown() {
	p.StopLoop()
	p.loopPlayer.Close()
}
