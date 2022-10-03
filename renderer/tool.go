package renderer

import (
	"github.com/HaBaLeS/arkhamassets"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"path"
)

func LoadImage(file string) *ebiten.Image {
	fh, err := arkhamassets.Data.Open(path.Join("data/", file))
	if err != nil {
		log.Panicln("Could not load image: %s: %v", file, err)
	}
	i, _, err := image.Decode(fh)
	if err != nil {
		log.Panicln("Could not load image: %s: %v", file, err)
	}
	return ebiten.NewImageFromImage(i)
}
