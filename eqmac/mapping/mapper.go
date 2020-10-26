package mapping

import (
	"github.com/google/uuid"
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
)

type Mapper interface {
	MapFixedBand(autoeq.FixedBandEQs, autoeq.EQMetadata) (eqmac.EQPreset, error)
}

// compile time interface implementation check
var _ Mapper = AutoEQMapper{}

type AutoEQMapper struct{}

func (m AutoEQMapper) MapFixedBand(fbeq autoeq.FixedBandEQs, meta autoeq.EQMetadata) (eqmac.EQPreset, error) {
	var preset eqmac.EQPreset
	preset.ID = uuid.New().String()
	preset.IsDefault = false
	preset.Name = meta.Name
	preset.Gains = eqmac.Gains{
		Global: meta.Global,
		Bands:  []float64{},
	}
	bands := []float64{}
	for _, band := range fbeq {
		bands = append(bands, band.Gain)
	}
	preset.Gains.Bands = bands
	return preset, nil
}
