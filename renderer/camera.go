package renderer

import (
	"ebiten2arkham/config"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Camera struct {
	//https://www.youtube.com/watch?v=LfbqtmqxX04
	Xoff           float64
	Yoff           float64
	viewportHeight float64
	viewportWifth  float64
	Rotation       float64
	ZoomFactor     float64
	filter         ebiten.Filter

	world *ebiten.Image
}

func NewCamera(world *ebiten.Image) *Camera {
	c := &Camera{
		Xoff:       0,
		Yoff:       0,
		ZoomFactor: 1,
		world:      world,
	}
	return c
}

func (c *Camera) WorldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Xoff, -c.Yoff)
	// We want to scale and rotate around center of image / screen
	x, y := c.viewportCenter()
	m.Translate(-x, -y)
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(x, y)
	return m
}

func (c *Camera) viewportCenter() (float64, float64) {
	return config.Cfg.CameraWidth / 2, config.Cfg.CameraHeight / 2
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.WorldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Xoff = 0
	c.Yoff = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}

func (c *Camera) Draw(screen *ebiten.Image) {
	screen.DrawImage(c.world, &ebiten.DrawImageOptions{
		Filter: c.filter,
		GeoM:   c.WorldMatrix(),
	})
}
