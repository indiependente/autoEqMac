//go:generate mockgen -package autoeq -source=mdparser.go -destination mdparser_mock.go

package autoeq

import (
	"bytes"
	"fmt"
)

// compile time interface implementation check
var _ MarkDownParser = MetadataParser{}

// EQMetadata represents EQ metadata.
type EQMetadata struct {
	ID     string
	Name   string
	Author string
	Link   string
	Global float64
}

// MarkDownParser defines the behaviour of a component capable of parsing raw bytes containing MarkDown text.
type MarkDownParser interface {
	ParseMetadata([]byte) ([]EQMetadata, error)
}

// MetadataParser is an implementation of a MarkDownParser that parses MarkDown coming from GitHub.
type MetadataParser struct {
	LinkPrefix        string
	FixedBandEQSuffix string
}

// ParseMetadata returns a slice of EQ metadata parsed from the input raw bytes.
// Returns an error if any.
func (p MetadataParser) ParseMetadata(data []byte) ([]EQMetadata, error) {
	var (
		eqMeta  []EQMetadata
		eqCount int
	)
	lines := bytes.Split(data, []byte("\n"))
	for _, l := range lines {
		if !bytes.HasPrefix(l, []byte("- [")) {
			continue
		}
		nameLinkAuth := bytes.Split(l, []byte("]("))
		name := bytes.TrimLeft(nameLinkAuth[0], "- [")
		linkAuth := bytes.Split(nameLinkAuth[1], []byte(") by "))
		link := p.LinkPrefix + string(bytes.TrimLeft(linkAuth[0], ".")) + buildLink(linkAuth[0]) + p.FixedBandEQSuffix
		eqMeta = append(eqMeta, EQMetadata{
			ID:     fmt.Sprintf("%d", eqCount),
			Name:   string(name),
			Author: string(linkAuth[1]),
			Link:   link,
		})
		eqCount++
	}
	return eqMeta, nil
}

func buildLink(partialLink []byte) string {
	return string(partialLink[bytes.LastIndex(partialLink, []byte("/")):])
}
