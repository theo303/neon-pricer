package svg

import "math"

// FormType is an enum to define types of forms.
type FormType string

const (
	RectangleType FormType = "rect"
	CircleType    FormType = "circle"
)

// Form defines a svg object for which a perimeter can be calculated.
type Form interface {
	GetPerimeter() float64
}

type Rectangle struct {
	Height, Width float64
	X, Y          float64
}

func (r Rectangle) GetPerimeter() float64 {
	return r.Height*2 + r.Width*2
}

type Circle struct {
	R    float64
	X, Y float64
}

func (c Circle) GetPerimeter() float64 {
	return 2 * math.Pi * c.R
}
