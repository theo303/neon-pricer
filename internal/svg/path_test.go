package svg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Path_Length(t *testing.T) {
	tests := map[string]struct {
		path Path
		want float64
	}{
		"1": {
			path: Path{
				Command:    'M',
				Parameters: []float64{10, 10},
				Next: &Path{
					Command:    'h',
					Parameters: []float64{10},
				},
			},
			want: 10,
		},
		"2": {
			path: Path{
				Command:    'M',
				Parameters: []float64{10, 10},
				Next: &Path{
					Command:    'm',
					Parameters: []float64{10, 10},
					Next: &Path{
						Command:    'H',
						Parameters: []float64{-20},
					},
				},
			},
			want: 40,
		},
		"3": {
			path: Path{
				Command:    'V',
				Parameters: []float64{20},
				Next: &Path{
					Command:    'v',
					Parameters: []float64{5},
				},
			},
			want: 25,
		},
		"4": {
			path: Path{
				Command:    'L',
				Parameters: []float64{3, -4},
				Next: &Path{
					Command:    'l',
					Parameters: []float64{5, 12},
				},
			},
			want: 18,
		},
		"5": {
			path: Path{
				Command:    'l',
				Parameters: []float64{3, 4, 4, 3},
			},
			want: 10,
		},
		"6": {
			path: Path{
				Command:    'L',
				Parameters: []float64{3, 4, 0, 8},
			},
			want: 10,
		},
		"7": {
			path: Path{
				Command:    'C',
				Parameters: []float64{110, 150, 25, 190, 210, 250, 210, 30},
			},
			want: 272.87,
		},
		"8": {
			path: Path{
				Command:    'Q',
				Parameters: []float64{220, 60, 20, 110, 70, 250},
			},
			want: 281.95,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tt.path.Length()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
			want:  point{x: 128.125, y: 187.5},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, splitBezier(tt.ratio, tt.points)[0])
		})
	}
}

func Test_parsePathCommand(t *testing.T) {
	tests := map[string]struct {
		pathString string
		want       Path
	}{
		"1 command with 2 parameters": {
			pathString: "M1479,815.22",
			want: Path{
				Command:    'M',
				Parameters: []float64{1479, 815.22},
			},
		},
		"3 commands": {
			pathString: "M715,371.73h3.29c26.16,3.52,97.36,16.63,161,75.59,73.24,67.81,88,151.38,91.47,175.83",
			want: Path{
				Command:    'M',
				Parameters: []float64{715, 371.73},
				Next: &Path{
					Command:    'h',
					Parameters: []float64{3.29},
					Next: &Path{
						Command:    'c',
						Parameters: []float64{26.16, 3.52, 97.36, 16.63, 161, 75.59, 73.24, 67.81, 88, 151.38, 91.47, 175.83},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := parsePathCommand(tt.pathString)
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func Test_parseParam(t *testing.T) {
	tests := map[string]struct {
		str  string
		want []float64
	}{
		"1 number": {
			str:  "815.22",
			want: []float64{815.22},
		},
		"1 negative number": {
			str:  "-815.22",
			want: []float64{-815.22},
		},
		"2 numbers separated by comma": {
			str:  "815.22,1429",
			want: []float64{815.22, 1429},
		},
		"2 numbers separated by minus": {
			str:  "815.22-1429",
			want: []float64{815.22, -1429},
		},
		"2 numbers with leading dots": {
			str:  ".98.89",
			want: []float64{.98, .89},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := parseParam(tt.str)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
