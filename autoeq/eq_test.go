package autoeq

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	athM50X = []byte(`[{"gains":{"global":-4.7,"bands":[-1.4,-0.8,-4.8,0.9,1.4,-1,-0.5,-0.2,4.5,-5.3]},"id":"ddbdb3ae-3556-4138-b829-579d5369f24d","isDefault":false,"name":"Audio-Technica ATH-M50x"}]`)
)

func TestEQHTTPGetter_GetEQ(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(athM50X)),
	}
	tests := []struct {
		name              string
		setupExpectations func(doer *MockDoer)
		meta              EQMetadata
		want              []byte
		wantErr           bool
	}{
		{
			name: "Happy path",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			meta: EQMetadata{
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
			name: "Sad path",
			setupExpectations: func(doer *MockDoer) {
				doer.EXPECT().Do(gomock.Any()).Return(nil, errors.New("connection error"))
			},
			meta: EQMetadata{
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
