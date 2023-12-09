package svg

import (
	"fmt"
	"strconv"

	"github.com/JoshVarga/svgparser"
)

type Rectangle struct {
	point
	height, width float64
}

func (r Rectangle) Length() (float64, error) {
	return r.height*2 + r.width*2, nil
}

func (r Rectangle) Size() (Size, error) {
	return Size{
		height: r.height,
		width:  r.width,
	}, nil
}

func parseRectangle(element svgparser.Element) (Rectangle, error) {
	height, err := strconv.ParseFloat(element.Attributes["height"], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("parsing height: %w", err)
	}
	width, err := strconv.ParseFloat(element.Attributes["width"], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("parsing witdh: %w", err)
	}
	var x, y float64
	_, ok := element.Attributes["x"]
	if !ok {
		x = 0
	} else {
		x, err = strconv.ParseFloat(element.Attributes["x"], 64)
		if err != nil {
			return Rectangle{}, fmt.Errorf("parsing x: %w", err)
		}
	}
	_, ok = element.Attributes["y"]
	if !ok {
		y = 0
	} else {
		y, err = strconv.ParseFloat(element.Attributes["y"], 64)
		if err != nil {
			return Rectangle{}, fmt.Errorf("parsing y: %w", err)
		}
	}

	return Rectangle{
		height: height,
		width:  width,
		point: point{
			x: x,
			y: y,
		},
	}, nil
}
