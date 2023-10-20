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

const acceptanceLimit = 0.01

type Path struct {
	Command    rune
	Parameters []float64
	Next       *Path
}

func (p Path) Length() (float64, error) {
	return p.length(point{})
}

func (p Path) length(lastPos point) (float64, error) {
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
		if len(p.Parameters)%2 != 0 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command L", len(p.Parameters))
		}
		for i := 0; i < len(p.Parameters); i += 2 {
			lx := lastPos.x - p.Parameters[i]
			ly := lastPos.y - p.Parameters[i+1]
			length += lengthLine(lx, ly)
			lastPos.x = p.Parameters[i]
			lastPos.y = p.Parameters[i+1]
		}
	case 'l':
		if len(p.Parameters)%2 != 0 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command l", len(p.Parameters))
		}
		for i := 0; i < len(p.Parameters); i += 2 {
			length += lengthLine(p.Parameters[i], p.Parameters[i+1])
			lastPos.x += p.Parameters[i]
			lastPos.y += p.Parameters[i+1]
		}
	case 'C':
		if len(p.Parameters)%6 != 0 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command C", len(p.Parameters))
		}
		for i := 0; i < len(p.Parameters); i += 6 {
			points := []point{
				lastPos,
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: p.Parameters[i+2], y: p.Parameters[i+3]},
				{x: p.Parameters[i+4], y: p.Parameters[i+5]},
			}
			length += lengthBezier(points)
			lastPos = points[3]
		}
	case 'c':
		if len(p.Parameters)%6 != 0 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command c", len(p.Parameters))
		}
		for i := 0; i < len(p.Parameters); i += 6 {
			points := []point{
				lastPos,
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: p.Parameters[i+2], y: p.Parameters[i+3]},
				{x: lastPos.x + p.Parameters[i+4], y: lastPos.y + p.Parameters[i+5]},
			}
			length += lengthBezier(points)
			lastPos = points[3]
		}
	case 'Q':
		if len(p.Parameters)%4 != 0 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command Q", len(p.Parameters))
		}
		for i := 0; i < len(p.Parameters); i += 4 {
			points := []point{
				lastPos,
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: p.Parameters[i+2], y: p.Parameters[i+3]},
			}
			length += lengthBezier(points)
			lastPos = points[2]
		}
	case 'q':
		if len(p.Parameters)%4 != 0 {
			return 0, fmt.Errorf("invalid number of parameters (%d) for command q", len(p.Parameters))
		}
		for i := 0; i < len(p.Parameters); i += 4 {
			points := []point{
				lastPos,
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
			}
			length += lengthBezier(points)
			lastPos = points[2]
		}
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

func lengthLine(lx, ly float64) float64 {
	return math.Sqrt(lx*lx + ly*ly)
}

func lengthBezier(points []point) float64 {
	var length float64
	lastLength := -1.0
	step := 0.5
	segments := []point{points[0], points[len(points)-1]}
	var newSegments []point
	for math.Abs(lastLength-length) > acceptanceLimit {
		lastLength = length
		length = 0
		newSegments = make([]point, (len(segments)-1)*2+1)
		for i := range newSegments {
			if i%2 == 0 {
				newSegments[i] = segments[i/2]
			} else {
				newSegments[i] = splitBezier(step*float64(i), points)[0]
			}
		}
		step /= 2.0
		for i := 0; i < len(newSegments)-1; i++ {
			length += lengthLine(newSegments[i+1].x-newSegments[i].x, newSegments[i+1].y-newSegments[i].y)
		}
		segments = append([]point{}, newSegments...)
	}

	return math.Round(length*100) / 100
}

func splitBezier(t float64, points []point) []point {
	if len(points) == 1 {
		return points
	}
	newpoints := make([]point, len(points)-1)
	for i := 0; i < len(newpoints); i++ {
		newpoints[i] = point{
			x: (1-t)*points[i].x + t*points[i+1].x,
			y: (1-t)*points[i].y + t*points[i+1].y,
		}
	}
	return splitBezier(t, newpoints)
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
