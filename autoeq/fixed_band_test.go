package autoeq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToFixedBandEQs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		want    FixedBandEQs
		wantErr bool
	}{
		{
			name:    "Happy path",
			data:    []byte("Filter 1: ON PK Fc 31 Hz Gain 5.8 dB Q 1.41\n"),
			want:    FixedBandEQs{{Frequency: 31, Gain: 5.8, Q: 1.41}},
			wantErr: false,
		},
		{
			name:    "Sad path - Freq not int",
			data:    []byte("Filter 1: ON PK Fc AB Hz Gain 5.8 dB Q 1.41"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Sad path - Gain not float",
			data:    []byte("Filter 1: ON PK Fc 31 Hz Gain AB dB Q 1.41"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Sad path - Q not float",
			data:    []byte("Filter 1: ON PK Fc 31 Hz Gain 5.8 dB Q AB"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFixedBandEQs(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
