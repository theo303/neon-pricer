package usecases

import (
	"testing"
	"theo303/neon-pricer/conf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getSiliconePricing(t *testing.T) {
	tests := map[string]struct {
		pricings []conf.Silicone
		id       string
		want     float64
		wantErr  bool
	}{
		"silicone with size": {
			pricings: []conf.Silicone{
				{
					SizeMm:        10,
					PricePerMeter: 5.3,
				},
			},
			id:   "10Mm",
			want: 5.3,
		},
		"standard 12mm silicone": {
			pricings: []conf.Silicone{
				{
					SizeMm:        12,
					PricePerMeter: 2.6,
				},
			},
			id:   "RGB",
			want: 2.6,
		},
		"invalid silicone": {
			id:      "78",
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			size, err := getSiliconeSize(tt.id)
			require.NoError(t, err)
			got, err := getSiliconePricing(tt.pricings, size)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
