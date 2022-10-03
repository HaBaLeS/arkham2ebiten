package input

import "github.com/hajimehoshi/ebiten/v2"

type Stroke struct {
	startx int
	starty int
	endx   int
	endy   int
	Active bool
}

func (s *Stroke) Start() {
	s.endx = 0
	s.starty = 0
	s.startx, s.starty = ebiten.CursorPosition()
	s.Active = true
}

func (s *Stroke) Stop() {
	s.endx, s.endy = ebiten.CursorPosition()
	s.Active = false
}

func (s *Stroke) Delta() (int, int) {
	s.endx, s.endy = ebiten.CursorPosition()
	dx := s.startx - s.endx
	dy := s.starty - s.endy
	s.startx = s.endx
	s.starty = s.endy
	return dx, dy
}
