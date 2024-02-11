package svg

import (
	"encoding/hex"
	"fmt"
	"io"
	"regexp"

	"github.com/JoshVarga/svgparser"
)

var regexpHexCode = regexp.MustCompile(`_[xX]([0-9a-fA-F]+)_`)

const decoupe = "DECOUPE"

// FormType is an enum to define types of forms.
type FormType string

const (
	RectangleType FormType = "rect"
	CircleType    FormType = "circle"
	PathType      FormType = "path"
	LineType      FormType = "line"
)

type point struct {
	x, y float64
}

// Form defines a svg object that can be measured, sized.
type Form interface {
	Length() (float64, error)
	Bounds() (Bounds, error)
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
	case string(LineType):
		line, err := parseLine(*element)
		if err != nil {
			return nil, fmt.Errorf("parsing line: %w", err)
		}
		forms = append(forms, line)
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

	formsGroups := make(map[string][]Form)
	for _, child := range element.Children {
		if child.Name != "g" {
			continue
		}
		if groupID != "" && child.Attributes["id"] != groupID {
			continue
		}

		groupID, err := sanitizeGroupID(child.Attributes["id"])
		if err != nil {
			return nil, fmt.Errorf("sanitizing group id %s: %w", child.Attributes["id"], err)
		}

		formsGroups[groupID], err = parseForms(child)
		if err != nil {
			return nil, fmt.Errorf("parsing group of forms %s: %w", child.Attributes["id"], err)
		}
	}

	return formsGroups, nil
}

func sanitizeGroupID(groupID string) (string, error) {
	loc := regexpHexCode.FindStringSubmatchIndex(groupID)
	for loc != nil {
		decodedStr, err := hex.DecodeString(groupID[loc[0]+2 : loc[1]-1])
		if err != nil {
			return "", fmt.Errorf("decoding %s: %w", groupID[loc[0]:loc[1]], err)
		}
		groupID = groupID[:loc[0]] + string(decodedStr) + groupID[loc[1]:]
		loc = regexpHexCode.FindStringIndex(groupID)
	}
	return groupID, nil
}
