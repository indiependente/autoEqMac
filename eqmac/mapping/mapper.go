package mapping

import (
	"github.com/google/uuid"
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
)

// Mapper defines the behaviour of a component capable of mapping AutoEQ data into EqMac format.
type Mapper interface {
	MapFixedBand(autoeq.FixedBandEQs, autoeq.EQMetadata) (eqmac.EQPreset, error)
}

// compile time interface implementation check
var _ Mapper = AutoEQMapper{}

// AutoEQMapper is an implementation of the Mapper interface.
// It maps autoeq FixedBandEQs into eqmac Presets.
type AutoEQMapper struct{}

// MapFixedBand maps fixed bands EQ data into an EqMac preset.
// Returns an error if any.
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
