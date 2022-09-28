package autoeq

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	athM50X       = []byte(`[{"gains":{"global":-4.7,"bands":[-1.4,-0.8,-4.8,0.9,1.4,-1,-0.5,-0.2,4.5,-5.3]},"id":"ddbdb3ae-3556-4138-b829-579d5369f24d","isDefault":false,"name":"Audio-Technica ATH-M50x"}]`)
	athM50XReadme = []byte(`# Audio-Technica ATH-M50x
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
)

func TestEQHTTPGetter_GetEQ(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(athM50X)),
	}
	tests := []struct {
		name              string
		setupExpectations func(doer *MockDoer)
		meta              *EQMetadata
		want              []byte
		wantErr           bool
	}{
		{
			name: "Happy path",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			meta: &EQMetadata{
				ID:     "123456789",
				Name:   "Audio-Technica ATH-M50x",
				Author: "jaakkopasanen",
				Link:   "https://github.com/jaakkopasanen/AutoEq/tree/master/results/oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x",
				Global: -6.4,
			},
			want:    athM50X,
			wantErr: false,
		},
		{
			name: "Sad path - connection error",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(nil, errors.New("connection error"))
			},
			meta: &EQMetadata{
				ID:     "123456789",
				Name:   "Audio-Technica ATH-M50x",
				Author: "jaakkopasanen",
				Link:   "https://github.com/jaakkopasanen/AutoEq/tree/master/results/oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x",
				Global: -6.4,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			doer := NewMockDoer(ctrl)
			tt.setupExpectations(doer)

			g := EQHTTPGetter{
				Client: doer,
			}
			got, err := g.GetEQ(tt.meta)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEQHTTPGetter_GetFixedBandGlobalPreamp(t *testing.T) { //nolint:funlen
	t.Parallel()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(athM50XReadme)),
	}
	tests := []struct {
		name              string
		setupExpectations func(doer *MockDoer)
		meta              *EQMetadata
		want              float64
		wantErr           bool
	}{
		{
			name: "Happy path",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			meta: &EQMetadata{
				ID:     "123456789",
				Name:   "Audio-Technica ATH-M50x",
				Author: "jaakkopasanen",
				Link:   "https://github.com/jaakkopasanen/AutoEq/tree/master/results/oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x",
			},
			want:    -6.4,
			wantErr: false,
		},
		{
			name: "Sad path - connection error",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bytes.ReplaceAll(athM50XReadme, []byte("dB"), []byte("db")))),
				}, nil)
			},
			meta: &EQMetadata{
				ID:     "123456789",
				Name:   "Audio-Technica ATH-M50x",
				Author: "jaakkopasanen",
				Link:   "https://github.com/jaakkopasanen/AutoEq/tree/master/results/oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Sad path - bad data error",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(nil, errors.New("connection error"))
			},
			meta: &EQMetadata{
				ID:     "123456789",
				Name:   "Audio-Technica ATH-M50x",
				Author: "jaakkopasanen",
				Link:   "https://github.com/jaakkopasanen/AutoEq/tree/master/results/oratory1990/harman_over-ear_2018/Audio-Technica%20ATH-M50x",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			doer := NewMockDoer(ctrl)
			tt.setupExpectations(doer)

			g := EQHTTPGetter{
				Client: doer,
			}
			got, err := g.GetFixedBandGlobalPreamp(tt.meta)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
