package gui

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
	"os"
)

var normalFnt font.Face

func InitFonts() {

	fnt, err := os.ReadFile("../data/font/Teutonic.ttf")
	if err != nil {
		log.Panicf("Could not Read font %v", err)
	}

	tt, err := opentype.Parse(fnt)
	if err != nil {
		log.Panicf("Could not Read font %v", err)
	}

	const dpi = 72
	normalFnt, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

}
