package svg

import (
	"errors"
	"math"
)

// https://stackoverflow.com/questions/9017100/calculate-center-of-svg-arc
// https://www.w3.org/TR/SVG/implnote.html#ArcImplementationNotes
func arcToCenterParam(start, end point, rx, ry, rot float64, fA, fS bool) (point, error) {
	if rx == 0 || ry == 0 {
		return point{}, errors.New("rx or ry cannot be equal to 0")
	}
	rx = math.Abs(rx)
	ry = math.Abs(ry)

	phi := rot * math.Pi / 180.0

	cosPhi := math.Cos(phi)
	sinPhi := math.Sin(phi)

	halfDiffX := (start.x - end.x) / 2.0
	halfDiffY := (start.y - end.y) / 2.0

	x1 := cosPhi*halfDiffX + sinPhi*halfDiffY
	y1 := -sinPhi*halfDiffX + cosPhi*halfDiffY

	rxy1 := rx * y1
	ryx1 := ry * x1

	coef := math.Sqrt((rx*rx*ry*ry - rxy1*rxy1 - ryx1*ryx1) / (rxy1*rxy1 + ryx1*ryx1))
	if fA == fS {
		coef = -coef
	}
	cx := coef * rxy1 / ry
	cy := -coef * ryx1 / rx

	halfSumX := (start.x + end.x) / 2
	halfSumY := (start.y + end.y) / 2
	center := point{
		x: cosPhi*cx - sinPhi*cy + halfSumX,
		y: sinPhi*cx + cosPhi*cy + halfSumY,
	}

	// x1 :=

	return center, nil
}
