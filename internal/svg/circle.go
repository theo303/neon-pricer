package svg

import (
	"fmt"
	"math"
	"strconv"

	"github.com/JoshVarga/svgparser"
)

type Circle struct {
	point
	r float64
}

func (c Circle) Length() (float64, error) {
	return 2 * math.Pi * c.r, nil
}

func (c Circle) Size() (Size, error) {
	return Size{
		width:  c.r,
		height: c.r,
	}, nil
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
		r: r,
		point: point{
			x: x,
			y: y,
		},
	}, nil
}
