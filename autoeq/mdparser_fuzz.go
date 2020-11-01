// +build gofuzz

package autoeq

const (
	autoEQResults   = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results`
	fixedBandSuffix = `%20FixedBandEQ.txt`
)

// TODO: move into fuzz branch
// not working right now due to: https://github.com/dvyukov/go-fuzz/issues/294

func Fuzz(in []byte) int {
	p := MetadataParser{
		LinkPrefix:        autoEQResults,
		FixedBandEQSuffix: fixedBandSuffix,
	}
	_, err := p.ParseMetadata(in)
	if err == nil {
		return 1 // interesting for fuzz
	}
	return 0 // normal for fuzz
}
