package svg

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/JoshVarga/svgparser"
)

var pathRegex = regexp.MustCompile(`[a-zA-Z][\d\-.,]*`)

var paramRegex = regexp.MustCompile(`-?[\d]*\.?[\d]+`)

const bezierStep = 100
const ellipticArcAngleStep = math.Pi / 100000

var nbOfParams = map[rune]int{
	'M': 2,
	'H': 1,
	'V': 1,
	'L': 2,
	'C': 6,
	'Q': 4,
	'S': 4,
	'T': 2,
	'A': 7,
	'Z': 0,
}

type Path struct {
	Command    rune
	Parameters []float64
	Next       *Path
}

func (p Path) checkNumberOfParams() bool {
	c := unicode.ToUpper(p.Command)
	n, ok := nbOfParams[c]
	if !ok {
		return false
	}
	if c == 'M' || c == 'V' || c == 'H' || c == 'Z' {
		return n == len(p.Parameters)
	}
	return len(p.Parameters)%n == 0
}

func (p Path) Length() (float64, error) {
	return p.length(point{}, point{}, point{})
}

func (p Path) length(firstPos, lastPos, lastCtrl point) (float64, error) {
	var length float64

	if !p.checkNumberOfParams() {
		return 0, fmt.Errorf("invalid number of parameters (%d) for command %c", len(p.Parameters), p.Command)
	}
	switch p.Command {
	case 'M':
		lastPos.x = p.Parameters[0]
		lastPos.y = p.Parameters[1]
		firstPos = lastPos
	case 'm':
		lastPos.x += p.Parameters[0]
		lastPos.y += p.Parameters[1]
		firstPos = lastPos
	case 'H':
		length = math.Abs(lastPos.x - p.Parameters[0])
		lastPos.x = p.Parameters[0]
	case 'h':
		length = math.Abs(p.Parameters[0])
		lastPos.x += p.Parameters[0]
	case 'V':
		length = math.Abs(lastPos.y - p.Parameters[0])
		lastPos.y = p.Parameters[0]
	case 'v':
		length = math.Abs(p.Parameters[0])
		lastPos.y += p.Parameters[0]
	case 'L':
		for i := 0; i < len(p.Parameters); i += 2 {
			lx := lastPos.x - p.Parameters[i]
			ly := lastPos.y - p.Parameters[i+1]
			length += lengthLine(lx, ly)
			lastPos.x = p.Parameters[i]
			lastPos.y = p.Parameters[i+1]
		}
	case 'l':
		for i := 0; i < len(p.Parameters); i += 2 {
			length += lengthLine(p.Parameters[i], p.Parameters[i+1])
			lastPos.x += p.Parameters[i]
			lastPos.y += p.Parameters[i+1]
		}
	case 'C':
		for i := 0; i < len(p.Parameters); i += 6 {
			points := []point{
				lastPos,
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: p.Parameters[i+2], y: p.Parameters[i+3]},
				{x: p.Parameters[i+4], y: p.Parameters[i+5]},
			}
			length += lengthBezier(points)
			lastPos = points[3]
			lastCtrl = points[2]
		}
	case 'c':
		for i := 0; i < len(p.Parameters); i += 6 {
			points := []point{
				lastPos,
				{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
				{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
				{x: lastPos.x + p.Parameters[i+4], y: lastPos.y + p.Parameters[i+5]},
			}
			length += lengthBezier(points)
			lastPos = points[3]
			lastCtrl = points[2]
		}
	case 'S':
		for i := 0; i < len(p.Parameters); i += 4 {
			points := []point{
				lastPos,
				reflectPoint(lastCtrl, lastPos),
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: p.Parameters[i+2], y: p.Parameters[i+3]},
			}
			length += lengthBezier(points)
			lastPos = points[3]
			lastCtrl = points[2]
		}
	case 's':
		for i := 0; i < len(p.Parameters); i += 4 {
			points := []point{
				lastPos,
				reflectPoint(lastCtrl, lastPos),
				{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
				{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
			}
			length += lengthBezier(points)
			lastPos = points[3]
			lastCtrl = points[2]
		}
	case 'Q':
		for i := 0; i < len(p.Parameters); i += 4 {
			points := []point{
				lastPos,
				{x: p.Parameters[i], y: p.Parameters[i+1]},
				{x: p.Parameters[i+2], y: p.Parameters[i+3]},
			}
			length += lengthBezier(points)
			lastPos = points[2]
			lastCtrl = points[1]
		}
	case 'q':
		for i := 0; i < len(p.Parameters); i += 4 {
			points := []point{
				lastPos,
				{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
				{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
			}
			length += lengthBezier(points)
			lastPos = points[2]
			lastCtrl = points[1]
		}
	case 'T':
		for i := 0; i < len(p.Parameters); i += 2 {
			points := []point{
				lastPos,
				reflectPoint(lastCtrl, lastPos),
				{x: p.Parameters[i], y: p.Parameters[i+1]},
			}
			length += lengthBezier(points)
			lastPos = points[2]
			lastCtrl = points[1]
		}
	case 't':
		for i := 0; i < len(p.Parameters); i += 2 {
			points := []point{
				lastPos,
				reflectPoint(lastCtrl, lastPos),
				{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
			}
			length += lengthBezier(points)
			lastPos = points[2]
			lastCtrl = points[1]
		}
	case 'A':
		for i := 0; i < len(p.Parameters); i += 7 {
			end := point{p.Parameters[i+5], p.Parameters[i+6]}
			arc, err := arcFromSVGParams(
				lastPos,
				end,
				p.Parameters[i], p.Parameters[i+1],
				p.Parameters[i+2],
				p.Parameters[i+3] == 1,
				p.Parameters[i+4] == 1,
			)
			if err != nil {
				return 0, fmt.Errorf("building arc: %w", err)
			}
			length += arc.length(ellipticArcAngleStep)
			lastPos = end
		}
	case 'a':
		for i := 0; i < len(p.Parameters); i += 7 {
			end := point{lastPos.x + p.Parameters[i+5], lastPos.y + p.Parameters[i+6]}
			arc, err := arcFromSVGParams(
				lastPos,
				end,
				p.Parameters[i], p.Parameters[i+1],
				p.Parameters[i+2],
				p.Parameters[i+3] == 1,
				p.Parameters[i+4] == 1,
			)
			if err != nil {
				return 0, fmt.Errorf("building arc: %w", err)
			}
			length += arc.length(ellipticArcAngleStep)
			lastPos = end
		}
	case 'Z', 'z':
		length += lengthLines([]point{lastPos, firstPos})
	default:
		fmt.Printf("Unrecognized path command %c\n", p.Command)
	}
	if p.Next != nil {
		l, err := p.Next.length(firstPos, lastPos, lastCtrl)
		if err != nil {
			return 0, err
		}
		length += l
	}
	return length, nil
}

func (p Path) Bounds() (Bounds, error) {
	var b Bounds
	var lastPos, lastCtrl point
	var stop bool
	for !stop {
		if !p.checkNumberOfParams() {
			return Bounds{}, fmt.Errorf("invalid number of parameters (%d) for command %c", len(p.Parameters), p.Command)
		}
		switch p.Command {
		case 'M':
			b.minX = p.Parameters[0]
			b.maxX = p.Parameters[0]
			b.minY = p.Parameters[1]
			b.maxY = p.Parameters[1]
			lastPos = point{p.Parameters[0], p.Parameters[1]}
		case 'm':
			b.minX = min(b.minX, lastPos.x+p.Parameters[0])
			b.maxX = max(b.maxX, lastPos.x+p.Parameters[0])
			b.minY = min(b.minY, lastPos.y+p.Parameters[1])
			b.maxY = max(b.maxY, lastPos.y+p.Parameters[1])
			lastPos = point{lastPos.x + p.Parameters[0], lastPos.y + p.Parameters[1]}
		case 'H':
			lastPos.x = p.Parameters[0]
			b = b.expandPoint(lastPos)
		case 'h':
			lastPos.x += p.Parameters[0]
			b = b.expandPoint(lastPos)
		case 'V':
			lastPos.y = p.Parameters[0]
			b = b.expandPoint(lastPos)
		case 'v':
			lastPos.y += p.Parameters[0]
			b = b.expandPoint(lastPos)
		case 'L':
			lastPos.x = p.Parameters[0]
			lastPos.y = p.Parameters[1]
			b = b.expandPoint(lastPos)
		case 'l':
			lastPos.x += p.Parameters[0]
			lastPos.y += p.Parameters[1]
			b = b.expandPoint(lastPos)
		case 'C':
			for i := 0; i < len(p.Parameters); i += nbOfParams[p.Command] {
				points := []point{
					lastPos,
					{x: p.Parameters[i], y: p.Parameters[i+1]},
					{x: p.Parameters[i+2], y: p.Parameters[i+3]},
					{x: p.Parameters[i+4], y: p.Parameters[i+5]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[3]
				lastCtrl = points[2]
			}
		case 'c':
			for i := 0; i < len(p.Parameters); i += 6 {
				points := []point{
					lastPos,
					{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
					{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
					{x: lastPos.x + p.Parameters[i+4], y: lastPos.y + p.Parameters[i+5]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[3]
				lastCtrl = points[2]
			}
		case 'S':
			for i := 0; i < len(p.Parameters); i += 4 {
				points := []point{
					lastPos,
					reflectPoint(lastCtrl, lastPos),
					{x: p.Parameters[i], y: p.Parameters[i+1]},
					{x: p.Parameters[i+2], y: p.Parameters[i+3]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[3]
				lastCtrl = points[2]
			}
		case 's':
			for i := 0; i < len(p.Parameters); i += 4 {
				points := []point{
					lastPos,
					reflectPoint(lastCtrl, lastPos),
					{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
					{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[3]
				lastCtrl = points[2]
			}
		case 'Q':
			for i := 0; i < len(p.Parameters); i += 4 {
				points := []point{
					lastPos,
					{x: p.Parameters[i], y: p.Parameters[i+1]},
					{x: p.Parameters[i+2], y: p.Parameters[i+3]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[2]
				lastCtrl = points[1]
			}
		case 'q':
			for i := 0; i < len(p.Parameters); i += 4 {
				points := []point{
					lastPos,
					{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
					{x: lastPos.x + p.Parameters[i+2], y: lastPos.y + p.Parameters[i+3]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[2]
				lastCtrl = points[1]
			}
		case 'T':
			for i := 0; i < len(p.Parameters); i += 2 {
				points := []point{
					lastPos,
					reflectPoint(lastCtrl, lastPos),
					{x: p.Parameters[i], y: p.Parameters[i+1]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[2]
				lastCtrl = points[1]
			}
		case 't':
			for i := 0; i < len(p.Parameters); i += 2 {
				points := []point{
					lastPos,
					reflectPoint(lastCtrl, lastPos),
					{x: lastPos.x + p.Parameters[i], y: lastPos.y + p.Parameters[i+1]},
				}
				b = b.Expand(boundsBezier(points))
				lastPos = points[2]
				lastCtrl = points[1]
			}
		case 'A':
			for i := 0; i < len(p.Parameters); i += 7 {
				end := point{p.Parameters[i+5], p.Parameters[i+6]}
				arc, err := arcFromSVGParams(
					lastPos,
					end,
					p.Parameters[i], p.Parameters[i+1],
					p.Parameters[i+2],
					p.Parameters[i+3] == 1,
					p.Parameters[i+4] == 1,
				)
				if err != nil {
					return Bounds{}, fmt.Errorf("building arc: %w", err)
				}
				b = b.Expand(arc.bounds(ellipticArcAngleStep))
				lastPos = end
			}
		case 'a':
			for i := 0; i < len(p.Parameters); i += 7 {
				end := point{lastPos.x + p.Parameters[i+5], lastPos.y + p.Parameters[i+6]}
				arc, err := arcFromSVGParams(
					lastPos,
					end,
					p.Parameters[i], p.Parameters[i+1],
					p.Parameters[i+2],
					p.Parameters[i+3] == 1,
					p.Parameters[i+4] == 1,
				)
				if err != nil {
					return Bounds{}, fmt.Errorf("building arc: %w", err)
				}
				b = b.Expand(arc.bounds(ellipticArcAngleStep))
				lastPos = end
			}
		}

		if p.Next != nil {
			p = *p.Next
		} else {
			stop = true
		}
	}
	return Bounds{
		minX: math.Round((b.minX)*100) / 100,
		maxX: math.Round((b.maxX)*100) / 100,
		minY: math.Round((b.minY)*100) / 100,
		maxY: math.Round((b.maxY)*100) / 100,
	}, nil
}

func lengthLine(lx, ly float64) float64 {
	return math.Sqrt(lx*lx + ly*ly)
}

func lengthLines(points []point) float64 {
	var length float64
	for i := 0; i < len(points)-1; i++ {
		length += lengthLine(points[i+1].x-points[i].x, points[i+1].y-points[i].y)
	}
	return length
}

// reflectPoint returns the reflection of p about m.
func reflectPoint(p, m point) point {
	return point{
		x: m.x*2 - p.x,
		y: m.y*2 - p.y,
	}
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
	pathString = strings.ReplaceAll(pathString, "\n", "")
	pathString = strings.ReplaceAll(pathString, " ", "")
	pathString = strings.ReplaceAll(pathString, "\t", "")
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
