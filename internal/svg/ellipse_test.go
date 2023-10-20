package svg

import (
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
		want  point
	}{
		"1": {
			start: point{10, 0},
			end:   point{0, 10},
			rx:    10,
			ry:    10,
			fA:    true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := arcToCenterParam(tt.start, tt.end, tt.rx, tt.ry, tt.rot, tt.fA, tt.fS)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
