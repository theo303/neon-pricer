package svg

import (
	"fmt"
	"math"
	"strconv"

	"github.com/JoshVarga/svgparser"
)

type Circle struct {
	R    float64
	X, Y float64
}

func (c Circle) Length() float64 {
	return 2 * math.Pi * c.R
}

func parseCircle(element svgparser.Element) (Circle, error) {
	r, err := strconv.ParseFloat(element.Attributes["r"], 64)
	if err != nil {
		return Circle{}, fmt.Errorf("parsing r: %w", err)
	}
	x, err := strconv.ParseFloat(element.Attributes["cx"], 64)
	if err != nil {
		return Circle{}, fmt.Errorf("parsing cx: %w", err)
	}
	y, err := strconv.ParseFloat(element.Attributes["cy"], 64)
	if err != nil {
		return Circle{}, fmt.Errorf("parsing cy: %w", err)
	}
	return Circle{
		R: r,
		X: x,
		Y: y,
	}, nil
}
