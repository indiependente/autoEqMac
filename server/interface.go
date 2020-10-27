package server

import (
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
)

type Server interface {
	ListEQsMetadata() ([]autoeq.EQMetadata, error)
	GetFixedBandEQPreset(id string) (eqmac.EQPreset, error)
	GetEQMetadataByName(name string) (autoeq.EQMetadata, error)
}
