package models

type Vector2D struct {
	x float64
	y float64
}


func NewVector2D(x, y float64) * Vector2D {
	return &Vector2D{ x, y, }
}

func (v * Vector2D) copy() *Vector2D {
	return &Vector2D {
		x: v.x, 
		y: v.y, 
	}
}

func (v* Vector2D) apply(other * Vector2D) {
	v.x *= other.x 
	v.y *= other.y
}