package config

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Allways available reasonalble config
// can be overwritten by file
var Cfg Configuration

const CFG_FILE_NAME = "arkham.cfg"

type Configuration struct {
	Fullscreen       bool
	Mute             bool
	CameraMoveFactor float64
	CameraWidth      float64
	CameraHeight     float64
	DefaultFilter    ebiten.Filter
}

func Init() {
	Cfg = Configuration{
		Fullscreen:       false,
		Mute:             false,
		CameraMoveFactor: 5,
		CameraWidth:      1920.0,
		CameraHeight:     1080.0,
		DefaultFilter:    ebiten.FilterNearest,
	}

	f, err := os.Open(CFG_FILE_NAME)
	if err != nil {
		log.Printf("Could not open Config. Using Defaults. Maybe file does not exist yet. %v", err)
		Save()
		return
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&Cfg)
	if err != nil {
		log.Panicf("%v", err)
	}

	//add any new params to the written file
	Save()
}

func Save() {
	f, err := os.Create(CFG_FILE_NAME)
	if err != nil {
		log.Fatalf("Error Saving config. %v", err)
	}
	defer f.Close()
	enc := yaml.NewEncoder(f)
	err = enc.Encode(Cfg)
	if err != nil {
		log.Panicf("%v", err)
	}
}
