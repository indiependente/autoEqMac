package server

import (
	"io"

	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
)

// Server defines the behaviour of a component capable of serving EQ related requests.
type Server interface {
	ListEQsMetadata() ([]autoeq.EQMetadata, error)
	GetFixedBandEQPreset(id string) (eqmac.EQPreset, error)
	GetEQMetadataByName(name string) (autoeq.EQMetadata, error)
	WritePreset(w io.Writer, p eqmac.EQPreset) error
}
