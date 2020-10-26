package autoeq

import (
	"bytes"
	"fmt"
)

// compile time interface implementation check
var _ MarkDownParser = MetadataParser{}

type EQMetadata struct {
	ID     string
	Name   string
	Author string
	Link   string
	Global float64
}

type MarkDownParser interface {
	ParseMetadata([]byte) ([]EQMetadata, error)
}

type MetadataParser struct {
	LinkPrefix        string
	FixedBandEQSuffix string
}

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
		name_link_auth := bytes.Split(l, []byte("]("))
		name := bytes.TrimLeft(name_link_auth[0], "- [")
		link_auth := bytes.Split(name_link_auth[1], []byte(") by "))
		link := p.LinkPrefix + string(bytes.TrimLeft(link_auth[0], ".")) + buildLink(link_auth[0]) + p.FixedBandEQSuffix
		eqMeta = append(eqMeta, EQMetadata{
			ID:     fmt.Sprintf("%d", eqCount),
			Name:   string(name),
			Author: string(link_auth[1]),
			Link:   link,
		})
		eqCount++
	}
	return eqMeta, nil
}

func buildLink(partialLink []byte) string {
	return string(partialLink[bytes.LastIndex(partialLink, []byte("/")):])
}
