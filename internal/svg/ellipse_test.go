package svg

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_arcToCenterParam(t *testing.T) {
	tests := map[string]struct {
		start point
		end   point
		rx    float64
		ry    float64
		rot   float64
		fA    bool
		fS    bool
		want  arc
	}{
		"1": {
			start: point{10, 0},
			end:   point{0, 10},
			rx:    10,
			ry:    10,
			fA:    true,
			want: arc{
				start:      point{10, 0},
				end:        point{0, 10},
				center:     point{0, 0},
				rx:         10,
				ry:         10,
				startAngle: 0,
				endAngle:   math.Pi / 2,
			},
		},
		"1.5": {
			start: point{10, 0},
			end:   point{0, 10},
			rx:    10,
			ry:    10,
			fA:    true,
			fS:    true,
			want: arc{
				start:      point{10, 0},
				end:        point{0, 10},
				center:     point{10, 10},
				rx:         10,
				ry:         10,
				startAngle: -math.Pi / 2,
				endAngle:   math.Pi,
				clockwise:  true,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := arcFromSVGParams(tt.start, tt.end, tt.rx, tt.ry, tt.rot, tt.fA, tt.fS)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_arc_point(t *testing.T) {
	tests := map[string]struct {
		arc  arc
		t    float64
		want point
	}{
		"1": {
			arc: arc{
				center: point{0, 0},
				rx:     10,
				ry:     10,
			},
			t:    0,
			want: point{10, 0},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arc.point(tt.t))
		})
	}
}

func Test_arc_length(t *testing.T) {
	tests := map[string]struct {
		arc  arc
		step float64
		want float64
	}{
		"1": {
			arc: arc{
				start:      point{10, 0},
				end:        point{0, 10},
				center:     point{0, 0},
				rx:         10,
				ry:         10,
				startAngle: 0,
				endAngle:   math.Pi / 2,
			},
			step: math.Pi / 10000,
			want: 47.123889610057255,
		},
		"2": {
			arc: arc{
				start:      point{10, 0},
				end:        point{0, 10},
				center:     point{0, 0},
				rx:         10,
				ry:         10,
				clockwise:  true,
				startAngle: 0,
				endAngle:   math.Pi / 2,
			},
			step: math.Pi / 10000,
			want: 15.704821610713136,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arc.length(tt.step))
		})
	}
}
