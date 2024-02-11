package domain

import "math"

type Size struct {
	Length   float64
	LengthPx float64
	Height   float64
	Width    float64
}

func Round(n float64) float64 {
	return math.Round(n*100) / 100
}
