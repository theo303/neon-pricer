package svg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/JoshVarga/svgparser"
)

var pathRegex = regexp.MustCompile(`[a-zA-Z][\d\-.,]*`)

var paramRegex = regexp.MustCompile(`-?[\d]*\.?[\d]+`)

type Path struct {
	Command    rune
	Parameters []float64
	Next       *Path
}

func (p Path) Length() float64 {
	var length float64
	if p.Next != nil {
		length = p.Next.Length()
	}
	switch p.Command {
	case 'M':
		return 0
	default:
		fmt.Printf("Unrecognized path command %c\n", p.Command)
	}
	return length
}

func newPath(str string) (*Path, error) {
	c := rune(str[0])
	params, err := parseParam(str[1:])
	if err != nil {
		return nil, fmt.Errorf("parsing params %s: %w", str[1:], err)
	}

	return &Path{
		Command:    c,
		Parameters: params,
	}, nil
}

func parsePath(element svgparser.Element) (Path, error) {
	pathString := element.Attributes["d"]
	path, err := parsePathCommand(pathString)
	if err != nil {
		return Path{}, err
	}
	if path == nil {
		path = &Path{}
	}
	return *path, nil
}

func parsePathCommand(pathString string) (*Path, error) {
	commandStr := pathRegex.FindString(pathString)
	if commandStr == "" {
		return nil, nil
	}
	path, err := newPath(commandStr)
	if err != nil {
		return nil, fmt.Errorf("parsing path %s: %w", commandStr, err)
	}
	nextString, ok := strings.CutPrefix(pathString, commandStr)
	if !ok {
		return nil, fmt.Errorf("prefix %s not found in %s", commandStr, pathString)
	}
	path.Next, err = parsePathCommand(nextString)
	return path, err
}

func parseParam(str string) ([]float64, error) {
	var params []float64
	for _, param := range paramRegex.FindAllString(str, -1) {
		n, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing param %s: %w", param, err)
		}
		params = append(params, n)
	}
	return params, nil
}
