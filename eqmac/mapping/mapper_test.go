package mapping

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
	"github.com/stretchr/testify/require"
)

func TestAutoEQMapper_MapFixedBand(t *testing.T) {
	t.Parallel()
	id := uuid.New().String()
	type args struct {
		fbeq autoeq.FixedBandEQs
		meta autoeq.EQMetadata
	}
	tests := []struct {
		name              string
		args              args
		setupExpectations func(gen *MockUUIDGenerator)
		want              eqmac.EQPreset
		wantErr           bool
	}{
		{
			name: "Happy path",
			args: args{
				fbeq: autoeq.FixedBandEQs{
					{
						Frequency: 31,
						Gain:      5.8,
						Q:         1.41,
					},
				},
				meta: autoeq.EQMetadata{
					ID:     "0",
					Name:   "ATH-M50x",
					Author: "mimmo",
					Link:   "https://link",
					Global: -6.4,
				},
			},
			setupExpectations: func(gen *MockUUIDGenerator) {
				gen.EXPECT().UUID().Return(id)
			},
			want: eqmac.EQPreset{
				Gains: eqmac.Gains{
					Global: -6.4,
					Bands:  []float64{5.8},
				},
				ID:        id,
				IsDefault: false,
				Name:      "ATH-M50x",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			gen := NewMockUUIDGenerator(ctrl)
			tt.setupExpectations(gen)
			m := AutoEQMapper{gen: gen}
			got, err := m.MapFixedBand(tt.args.fbeq, tt.args.meta)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
