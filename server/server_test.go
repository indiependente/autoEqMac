package server

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
	"github.com/stretchr/testify/require"
)

var (
	rawEqList = []byte(`# Index
This is a list of all equalization profiles. Target is in parentheses if there are results with multiple targets
from the same source.

- [Audio-Technica ATH-M50x](./oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x) by oratory1990
`)
	rawEqData     = []byte(`Filter 1: ON PK Fc 31 Hz Gain 5.8 dB Q 1.41`)
	rawGlobalData = []byte(`# Audio-Technica ATH-M50x
See [usage instructions](https://github.com/jaakkopasanen/AutoEq#usage) for more options and info.

### Parametric EQs
In case of using parametric equalizer, apply preamp of **-7.0dB** and build filters manually
with these parameters. The first 5 filters can be used independently.
When using independent subset of filters, apply preamp of **-7.0dB**.

| Type    | Fc       |    Q | Gain    |
|:--------|:---------|:-----|:--------|
| Peaking | 23 Hz    | 0.95 | 6.3 dB  |
| Peaking | 327 Hz   | 2.37 | 3.2 dB  |
| Peaking | 5826 Hz  | 5.21 | 6.7 dB  |
| Peaking | 18679 Hz | 0.06 | -2.4 dB |
| Peaking | 19122 Hz | 0.41 | -9.3 dB |
| Peaking | 167 Hz   | 2.59 | -2.0 dB |
| Peaking | 1397 Hz  | 0.62 | 1.2 dB  |
| Peaking | 2608 Hz  | 3.12 | -2.3 dB |
| Peaking | 7966 Hz  | 2.55 | -2.1 dB |
| Peaking | 9166 Hz  | 3.3  | 3.4 dB  |

### Fixed Band EQs
In case of using fixed band (also called graphic) equalizer, apply preamp of **-6.4dB**
(if available) and set gains manually with these parameters.

| Type    | Fc       |    Q | Gain     |
|:--------|:---------|:-----|:---------|
| Peaking | 31 Hz    | 1.41 | 5.8 dB   |
| Peaking | 62 Hz    | 1.41 | 0.9 dB   |
| Peaking | 125 Hz   | 1.41 | -1.7 dB  |
| Peaking | 250 Hz   | 1.41 | 1.7 dB   |
| Peaking | 500 Hz   | 1.41 | 0.9 dB   |
| Peaking | 1000 Hz  | 1.41 | 0.9 dB   |
| Peaking | 2000 Hz  | 1.41 | -0.7 dB  |
| Peaking | 4000 Hz  | 1.41 | -0.8 dB  |
| Peaking | 8000 Hz  | 1.41 | -0.6 dB  |
| Peaking | 16000 Hz | 1.41 | -13.5 dB |

### Graphs
![](./Audio-Technica%20ATH-M50x.png)`)
	id     = uuid.New().String()
	eqMeta = autoeq.EQMetadata{
		ID:     "0",
		Name:   "Audio-Technica ATH-M50x",
		Author: "oratory1990",
		Link:   "https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results/oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x/Audio-Technica%20ATH-M50x%20FixedBandEQ.txt",
		Global: 0,
	}
	eqPreset = eqmac.EQPreset{
		Gains: eqmac.Gains{
			Global: -6.4,
			Bands:  []float64{5.8},
		},
		ID:        id,
		IsDefault: false,
		Name:      "Audio-Technica ATH-M50x",
	}
)

func TestHTTPServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		setupExpectations func(doer *MockDoer)
		want              struct {
			meta   []autoeq.EQMetadata
			preset eqmac.EQPreset
		}
		wantErr bool
	}{
		{
			name: "Happy path",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(rawEqList)),
				}, nil)
				doer.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(rawEqData)),
				}, nil)
				doer.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(rawGlobalData)),
				}, nil)
			},
			want: struct {
				meta   []autoeq.EQMetadata
				preset eqmac.EQPreset
			}{
				meta: []autoeq.EQMetadata{
					eqMeta,
				},
				preset: eqPreset,
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			doer := NewMockDoer(ctrl)
			mdp := autoeq.NewMetadataParser()
			eqg := autoeq.EQHTTPGetter{Client: doer}
			mapp := mapping.NewAutoEQMapper(mapping.WrappedGenerator(func() string {
				return id
			}))
			tt.setupExpectations(doer)

			s := HTTPServer{
				client:   doer,
				mdparser: mdp,
				eqGetter: eqg,
				mapper:   mapp,
				eqMetas:  map[string]autoeq.EQMetadata{},
				eqNameID: map[string]string{},
			}
			got, err := s.ListEQsMetadata()
			require.NoError(t, err)
			require.Equal(t, tt.want.meta, got)
			gotMeta, err := s.GetEQMetadataByName(eqMeta.Name)
			require.NoError(t, err)
			require.Equal(t, tt.want.meta[0], gotMeta)
			gotPreset, err := s.GetFixedBandEQPreset(eqMeta.ID)
			require.NoError(t, err)
			require.Equal(t, tt.want.preset, gotPreset)
			err = s.WritePreset(io.Discard, gotPreset)
			require.NoError(t, err)
		})
	}
}
