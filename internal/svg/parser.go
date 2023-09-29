package svg

import (
	"fmt"
	"io"
	"strconv"

	"github.com/JoshVarga/svgparser"
)

// RetrieveForms retrieves a list of Forms from the svg source.
func RetrieveForms(source io.Reader) ([]Form, error) {
	svg, err := svgparser.Parse(source, true)
	if err != nil {
		return nil, fmt.Errorf("parsing svg file: %w", err)
	}

	return findForms(svg)
}

func findForms(element *svgparser.Element) ([]Form, error) {
	if element == nil {
		return nil, nil
	}

	var forms []Form
	switch element.Name {
	case string(RectangleType):
		rect, err := parseRectangle(*element)
		if err != nil {
			return nil, fmt.Errorf("parsing rectangle: %w", err)
		}
		forms = append(forms, rect)
	case string(CircleType):
		circle, err := parseCircle(*element)
		if err != nil {
			return nil, fmt.Errorf("parsing circle: %w", err)
		}
		forms = append(forms, circle)
	}

	for i, child := range element.Children {
		childForms, err := findForms(child)
		if err != nil {
			return nil, fmt.Errorf("searching forms in child %d: %w", i, err)
		}
		forms = append(forms, childForms...)
	}
	return forms, nil
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
	x, err := strconv.ParseFloat(element.Attributes["x"], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("parsing x: %w", err)
	}
	y, err := strconv.ParseFloat(element.Attributes["y"], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("parsing y: %w", err)
	}
	return Rectangle{
		Height: height,
		Width:  width,
		X:      x,
		Y:      y,
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
		R: r,
		X: x,
		Y: y,
	}, nil
}
