package svg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeGroupID(t *testing.T) {
	tests := map[string]struct {
		groupID string
		want    string
	}{
		"_x38_MM": {
			groupID: "_x38_MM",
			want:    "8MM",
		},
		"_x38_MM_X13_": {
			groupID: "_x38_MM_X67_",
			want:    "8MMg",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := sanitizeGroupID(tt.groupID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
