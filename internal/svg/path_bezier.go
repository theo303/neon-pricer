package svg

import (
	"math"
)

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

func lengthBezier(points []point) float64 {
	var length float64
	lastLength := -1.0
	step := 0.5
	segments := []point{points[0], points[len(points)-1]}
	var newSegments []point
	for math.Abs(lastLength-length) > 1.0/bezierPrecision {
		lastLength = length
		newSegments = make([]point, (len(segments)-1)*2+1)
		for i := range newSegments {
			if i%2 == 0 {
				newSegments[i] = segments[i/2]
			} else {
				newSegments[i] = splitBezier(step*float64(i), points)[0]
			}
		}
		step /= 2.0
		length = lengthLines(newSegments)
		segments = append([]point{}, newSegments...)
	}

	return math.Round(length*100) / 100
}

func bezierXsAndYs(points []point) [2][bezierPrecision]float64 {
	positions := [2][bezierPrecision]float64{
		{points[0].x},
		{points[0].y},
	}
	for i := 1; i < bezierPrecision; i++ {
		p := splitBezier(float64(i)/bezierPrecision, points)[0]
		positions[0][i] = p.x
		positions[1][i] = p.y
	}
	return positions
}
