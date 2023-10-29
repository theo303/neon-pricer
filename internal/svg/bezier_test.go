package svg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitBezier(t *testing.T) {
	tests := map[string]struct {
		ratio  float64
		points []point
		want   point
	}{
		"1": {
			points: []point{
				{x: 110, y: 150},
				{x: 25, y: 190},
				{x: 210, y: 250},
				{x: 210, y: 30},
			},
			ratio: 0.5,
			want:  point{x: 128.125, y: 187.5},
		},
		"2": {
			points: []point{
				{x: 110, y: 150},
				{x: 25, y: 190},
				{x: 210, y: 250},
				{x: 210, y: 30},
			},
			ratio: 0.25,
			want:  point{x: 89.765625, y: 179.0625},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, splitBezier(tt.ratio, tt.points)[0])
		})
	}
}
