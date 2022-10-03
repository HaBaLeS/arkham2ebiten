package math

//FIXME write a test for math package

type Collider interface {
	Contains(x, y float64) bool
}

type Rectangle struct {
	X, Y, W, H float64
	Collider
}

func (r *Rectangle) Contains(x, y float64) bool {
	if x > r.X && x < r.X+r.W && y > r.Y && y < r.Y+r.H {
		return true
	}
	return false
}
