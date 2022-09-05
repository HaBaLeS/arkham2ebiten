package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"os"
	"path"
)

func LoadImage(file string) *ebiten.Image {
	fh, err := os.Open(path.Join("../data/", file))
	if err != nil {
		log.Panicln("Could not load image: %s: %v", file, err)
	}
	i, _, err := image.Decode(fh)
	if err != nil {
		log.Panicln("Could not load image: %s: %v", file, err)
	}
	return ebiten.NewImageFromImage(i)
}
