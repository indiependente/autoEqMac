package autoeq

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataParser_ParseMetadata(t *testing.T) {
	t.Parallel()
	type fields struct {
		LinkPrefix        string
		FixedBandEQSuffix string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*EQMetadata
		wantErr bool
	}{
		{
			name: "Happy path",
			fields: fields{
				LinkPrefix:        eqResultsPrefix,
				FixedBandEQSuffix: fixedBandSuffix,
			},
			args: args{
				data: mustReadFixture(t, "testdata/autoeq_index_top.txt"),
			},
			want: []*EQMetadata{
				{
					ID:     "0",
					Name:   "1Custom SA02",
					Author: "Crinacle",
					Link:   eqResultsPrefix + "/crinacle/harman_in-ear_2019v2/1Custom%20SA02/1Custom%20SA02" + fixedBandSuffix,
					Global: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := MetadataParser{
				LinkPrefix:        tt.fields.LinkPrefix,
				FixedBandEQSuffix: tt.fields.FixedBandEQSuffix,
			}
			got, err := p.ParseMetadata(tt.args.data)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}

func mustReadFixture(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%q\n", err)
		t.Fail()
	}

	return data
}
