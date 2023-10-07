package svg

import (
	"fmt"
	"io"

	"github.com/JoshVarga/svgparser"
)

// FormType is an enum to define types of forms.
type FormType string

const (
	RectangleType FormType = "rect"
	CircleType    FormType = "circle"
	PathType      FormType = "path"
)

// Measurable defines a svg object for which a length can be calculated.
type Measurable interface {
	Length() float64
}

// RetrieveForms retrieves a list of Forms from the svg source.
func RetrieveForms(source io.Reader, groupID string) ([]Measurable, error) {
	svg, err := svgparser.Parse(source, true)
	if err != nil {
		return nil, fmt.Errorf("parsing svg file: %w", err)
	}

	return findForms(svg, groupID)
}

func findForms(element *svgparser.Element, groupID string) ([]Measurable, error) {
	if element == nil || (groupID != "" && element.Name == "g" && element.Attributes["id"] != groupID) {
		return nil, nil
	}

	var forms []Measurable
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
	case string(PathType):
		path, err := parsePath(*element)
		if err != nil {
			return nil, fmt.Errorf("parsing path: %w", err)
		}
		forms = append(forms, path)
	}

	for i, child := range element.Children {
		childForms, err := findForms(child, groupID)
		if err != nil {
			return nil, fmt.Errorf("searching forms in child %d: %w", i, err)
		}
		forms = append(forms, childForms...)
	}
	return forms, nil
}
