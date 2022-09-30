//go:generate mockgen -package autoeq -source=mdparser.go -destination mdparser_mock.go

package autoeq

import (
	"bytes"
	"fmt"
)

const (
	eqResultsPrefix = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results`
	fixedBandSuffix = `%20FixedBandEQ.txt`
)

// compile time interface implementation check.
var _ MarkDownParser = MetadataParser{}

// EQMetadata represents EQ metadata.
type EQMetadata struct {
	ID     string
	Name   string
	Author string
	Link   string
	// Deprecated: this field has been deprecated in favor of autoeq.FixedBandEQ.Preamp.
	Global float64
}

// MarkDownParser defines the behavior of a component capable of parsing raw bytes containing MarkDown text.
type MarkDownParser interface {
	ParseMetadata([]byte) ([]*EQMetadata, error)
}

// MetadataParser is an implementation of a MarkDownParser that parses MarkDown coming from GitHub.
type MetadataParser struct {
	LinkPrefix        string
	FixedBandEQSuffix string
}

// NewMetadataParser returns a MetadataParser with populated fields.
func NewMetadataParser() MetadataParser {
	return MetadataParser{
		LinkPrefix:        eqResultsPrefix,
		FixedBandEQSuffix: fixedBandSuffix,
	}
}

// ParseMetadata returns a slice of EQ metadata parsed from the input raw bytes.
// Returns an error if any.
func (p MetadataParser) ParseMetadata(data []byte) ([]*EQMetadata, error) {
	eqCount := 0
	lines := bytes.Split(data, []byte("\n"))
	eqMeta := make([]*EQMetadata, len(lines))

	for _, l := range lines {
		if !bytes.HasPrefix(l, []byte("- [")) {
			continue
		}
		nameLinkAuth := bytes.Split(l, []byte("]("))
		name := bytes.TrimLeft(nameLinkAuth[0], "- [")
		linkAuth := bytes.Split(nameLinkAuth[1], []byte(") by "))
		link := p.LinkPrefix + string(bytes.TrimLeft(linkAuth[0], ".")) + buildLink(linkAuth[0]) + p.FixedBandEQSuffix
		eqMeta[eqCount] = &EQMetadata{
			ID:     fmt.Sprintf("%d", eqCount),
			Name:   string(name),
			Author: string(linkAuth[1]),
			Link:   link,
		}
		eqCount++
	}

	return eqMeta[:eqCount], nil
}

func buildLink(partialLink []byte) string {
	return string(partialLink[bytes.LastIndex(partialLink, []byte("/")):])
}
