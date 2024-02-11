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

func (c Circle) Bounds() (Bounds, error) {
	return Bounds{
		minX: c.x - c.r,
		maxX: c.x + c.r,
		minY: c.y - c.r,
		maxY: c.y + c.r,
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
