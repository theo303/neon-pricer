package svg

import (
	"fmt"
	"math"
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

func (p Path) Length() (float64, error) {
	return p.length(position{})
}

func (p Path) length(lastPos position) (float64, error) {
	var length float64

	switch p.Command {
	case 'M':
		if len(p.Parameters) != 2 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command M", len(p.Parameters))
		}
		lastPos.x = p.Parameters[0]
		lastPos.y = p.Parameters[1]
	case 'm':
		if len(p.Parameters) != 2 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command m", len(p.Parameters))
		}
		lastPos.x += p.Parameters[0]
		lastPos.y += p.Parameters[1]
	case 'H':
		if len(p.Parameters) != 1 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command H", len(p.Parameters))
		}
		length = math.Abs(lastPos.x - p.Parameters[0])
		lastPos.x = p.Parameters[0]
	case 'h':
		if len(p.Parameters) != 1 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command h", len(p.Parameters))
		}
		length = p.Parameters[0]
		lastPos.x += p.Parameters[0]
	case 'V':
		if len(p.Parameters) != 1 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command V", len(p.Parameters))
		}
		length = math.Abs(lastPos.y - p.Parameters[0])
		lastPos.y = p.Parameters[0]
	case 'v':
		if len(p.Parameters) != 1 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command v", len(p.Parameters))
		}
		length = p.Parameters[0]
		lastPos.y += p.Parameters[0]
	case 'L':
		if len(p.Parameters) != 2 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command L", len(p.Parameters))
		}
		lx := lastPos.x - p.Parameters[0]
		ly := lastPos.y - p.Parameters[1]
		length = math.Sqrt(lx*lx + ly*ly)
		lastPos.x = p.Parameters[0]
		lastPos.y = p.Parameters[1]
	case 'l':
		if len(p.Parameters) != 2 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command l", len(p.Parameters))
		}
		length = math.Sqrt(p.Parameters[0]*p.Parameters[0] + p.Parameters[1]*p.Parameters[1])
		lastPos.x += p.Parameters[0]
		lastPos.y += p.Parameters[1]
	default:
		fmt.Printf("Unrecognized path command %c\n", p.Command)
	}
	if p.Next != nil {
		l, err := p.Next.length(lastPos)
		if err != nil {
			return 0, err
		}
		length += l
	}
	return length, nil
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
