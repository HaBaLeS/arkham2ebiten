package font

import (
	"github.com/HaBaLeS/arkhamassets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

var NormalFnt font.Face

func InitFonts() {

	fnt, err := arkhamassets.Data.ReadFile("data/font/Teutonic.ttf")
	if err != nil {
		log.Panicf("Could not Read font %v", err)
	}

	tt, err := opentype.Parse(fnt)
	if err != nil {
		log.Panicf("Could not Read font %v", err)
	}

	const dpi = 72
	NormalFnt, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

}
