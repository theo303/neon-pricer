package svg

import (
	"fmt"
	"io"

	"github.com/JoshVarga/svgparser"
)

const decoupe = "DECOUPE"

// FormType is an enum to define types of forms.
type FormType string

const (
	RectangleType FormType = "rect"
	CircleType    FormType = "circle"
	PathType      FormType = "path"
)

type point struct {
	x, y float64
}

// Form defines a svg object that can be measured, sized.
type Form interface {
	Length() (float64, error)
	Size() (Size, error)
}

type Size struct {
	height, width float64
}

// RetrieveForms retrieves a list of Forms from the svg source.
func RetrieveForms(source io.Reader, groupID string) (map[string][]Form, error) {
	svg, err := svgparser.Parse(source, true)
	if err != nil {
		return nil, fmt.Errorf("parsing svg file: %w", err)
	}

	return parseGroups(svg, groupID)
}

func parseForms(element *svgparser.Element) ([]Form, error) {
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
	case string(PathType):
		path, err := parsePath(*element)
		if err != nil {
			return nil, fmt.Errorf("parsing path: %w", err)
		}
		forms = append(forms, path)
	}

	for i, child := range element.Children {
		childForms, err := parseForms(child)
		if err != nil {
			return nil, fmt.Errorf("searching forms in child %d: %w", i, err)
		}
		forms = append(forms, childForms...)
	}

	return forms, nil
}

func parseGroups(element *svgparser.Element, groupID string) (map[string][]Form, error) {
	if element == nil ||
		(groupID != "" && element.Name == "g" && element.Attributes["id"] != groupID) ||
		(element.Name == "g" && element.Attributes["id"] == decoupe) {
		return nil, nil
	}

	var err error
	formsGroups := make(map[string][]Form)
	for _, child := range element.Children {
		if child.Name != "g" {
			continue
		}
		if groupID != "" && child.Attributes["id"] != groupID {
			continue
		}
		formsGroups[child.Attributes["id"]], err = parseForms(child)
		if err != nil {
			return nil, fmt.Errorf("parsing group of forms %s: %w", child.Attributes["id"], err)
		}
	}

	return formsGroups, nil
}
