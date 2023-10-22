package svg

import (
	"errors"
	"math"
)

type arc struct {
	// center of the ellipse.
	center point
	// start and end points.
	start, end point
	// rx and ry are the radii of the ellipse.
	rx, ry float64
	// phi is the angle from the x-axis of the coordonate system to the x-axis of the ellipse.
	phi float64
	// angles of the arc.
	startAngle, endAngle float64

	clockwise bool
}

// compute angle between two vectors.
func angle(ux, uy, vx, vy float64) float64 {
	dot := ux*vx + uy*vy
	mod := math.Sqrt((ux*ux + uy*uy) * (vx*vx + vy*vy))
	angle := math.Acos(dot / mod)
	if ux*vy-uy*vx < 0 {
		return -angle
	}
	return angle
}

// https://stackoverflow.com/questions/9017100/calculate-center-of-svg-arc
// https://www.w3.org/TR/SVG/implnote.html#ArcImplementationNotes
func arcFromSVGParams(start, end point, rx, ry, rot float64, fA, fS bool) (arc, error) {
	if rx == 0 || ry == 0 {
		return arc{}, errors.New("rx or ry cannot be equal to 0")
	}
	rx = math.Abs(rx)
	ry = math.Abs(ry)

	phi := rot * math.Pi / 180

	cosPhi := math.Cos(phi)
	sinPhi := math.Sin(phi)

	halfDiffX := (start.x - end.x) / 2
	halfDiffY := (start.y - end.y) / 2

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

	xcr1 := (x1 - cx) / rx
	ycr1 := (y1 - cy) / ry
	startAngle := angle(1, 0, xcr1, ycr1)

	xcr2 := (-x1 - cx) / rx
	ycr2 := (-y1 - cy) / rx
	endAngle := math.Abs(math.Mod(startAngle+angle(xcr1, ycr1, xcr2, ycr2), math.Pi*2))

	return arc{
		center:     center,
		start:      start,
		end:        end,
		rx:         rx,
		ry:         ry,
		phi:        phi,
		startAngle: startAngle,
		endAngle:   endAngle,
		clockwise:  fS,
	}, nil
}

func (a arc) point(t float64) point {
	cosPhi := math.Cos(a.phi)
	sinPhi := math.Sin(a.phi)
	rxCosT := a.rx * math.Cos(t)
	rySinT := a.ry * math.Sin(t)
	return point{
		x: cosPhi*rxCosT - sinPhi*rySinT + a.center.x,
		y: sinPhi*rxCosT + cosPhi*rySinT + a.center.y,
	}
}

func (a arc) length(step float64) float64 {
	deltaAngle := math.Abs(math.Mod(a.endAngle-a.startAngle, math.Pi*2))
	if !a.clockwise {
		deltaAngle = math.Pi*2 - deltaAngle
	}
	points := []point{a.start}
	for t := step; t < deltaAngle; t += step {
		points = append(points, a.point(a.startAngle+t))
	}

	length := lengthLines(points)

	return length
}
