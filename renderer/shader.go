package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var shaderSource = `
package main

var Time float
var Cursor vec2
var ScreenSize vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {

	de(vec3(1,2,3))

	px := vec4(0.5,0.5,0.5,0.7)
	return px
}


func  de(p vec3 ) float{
	return 2.0
}

`

func GetShader() *ebiten.Shader {

	sh, err := ebiten.NewShader([]byte(shaderSource))
	if err != nil {
		log.Panicf("error creating shader %v", err)
	}

	return sh
}
