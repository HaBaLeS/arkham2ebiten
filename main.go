package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type Game struct {
	card *ebiten.Image
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(400, 600)

	screen.DrawImage(g.card, op)

}

func (g *Game) Layout(ow, oh int) (sh, sw int) {
	return ow, oh
}

func main() {
	//ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Arkham-go")
	ebiten.SetFullscreen(true)
	game := &Game{}

	f, err := os.Open("/home/falko/projekte/arkham-go/leech-img/" + "05186.png")
	if err != nil {
		panic(err)
	}
	i, n, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	log.Printf("Loaded Image with format %s", n)

	game.card = ebiten.NewImageFromImage(i)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
