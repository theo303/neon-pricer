package svg

import (
	"fmt"
	"math"
	"strconv"

	"github.com/JoshVarga/svgparser"
)

type Line struct {
	p1, p2 point
}

func (l Line) Length() (float64, error) {
	lx := l.p1.x - l.p2.x
	ly := l.p1.y - l.p2.y
	return math.Sqrt(lx*lx + ly*ly), nil
}

func (l Line) Size() (Size, error) {
	return Size{
		width:  max(l.p1.x, l.p2.x) - min(l.p1.x, l.p2.x),
		height: max(l.p1.y, l.p2.y) - min(l.p1.y, l.p2.y),
	}, nil
}

func parseLine(element svgparser.Element) (Line, error) {
	x1, err := strconv.ParseFloat(element.Attributes["x1"], 64)
	if err != nil {
		return Line{}, fmt.Errorf("parsing x1: %w", err)
	}
	y1, err := strconv.ParseFloat(element.Attributes["y1"], 64)
	if err != nil {
		return Line{}, fmt.Errorf("parsing y1: %w", err)
	}
	x2, err := strconv.ParseFloat(element.Attributes["x2"], 64)
	if err != nil {
		return Line{}, fmt.Errorf("parsing x2: %w", err)
	}
	y2, err := strconv.ParseFloat(element.Attributes["y2"], 64)
	if err != nil {
		return Line{}, fmt.Errorf("parsing y2: %w", err)
	}
	return Line{
		p1: point{x1, y1},
		p2: point{x2, y2},
	}, nil
}
