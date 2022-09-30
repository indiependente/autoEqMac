package autoeq

import (
	"fmt"
	"strconv"
	"strings"
)

const fixedBandFields = 12

// FixedBandFilter represents a single EQ band.
type FixedBandFilter struct {
	Frequency int     // Hz
	Gain      float64 // dB
	Q         float64 // fixed
}

// FixedBandEQ represents a simple fixed bands EQ.
type FixedBandEQ struct {
	Filters []*FixedBandFilter
	Preamp  float64
}

// ToFixedBandEQs transforms the raw EQ data into a fixed band EQ.
// Returns an error if any.
func ToFixedBandEQs(data []byte) (*FixedBandEQ, error) {
	rows := strings.Split(string(data), "\n")
	fbEQ := &FixedBandEQ{
		Filters: make([]*FixedBandFilter, len(rows)),
	}

	i := 0
	for _, row := range rows {
		if row == "" {
			continue
		}
		if strings.HasPrefix(row, "Preamp") {
			preamp, err := strconv.ParseFloat(strings.TrimSpace(strings.Fields(row)[1]), bitSize)
			if err != nil {
				return nil, err
			}
			fbEQ.Preamp = preamp

			continue
		}
		eqFields := strings.Fields(row)
		if len(eqFields) < fixedBandFields {
			return nil, fmt.Errorf("could not parse : %s", row)
		}
		freq, err := strconv.Atoi(strings.TrimSpace(eqFields[5]))
		if err != nil {
			return nil, fmt.Errorf("could not parse frequency: %w", err)
		}
		gain, err := strconv.ParseFloat(strings.TrimSpace(eqFields[8]), bitSize)
		if err != nil {
			return nil, fmt.Errorf("could not parse gain: %w", err)
		}
		q, err := strconv.ParseFloat(strings.TrimSpace(eqFields[11]), bitSize)
		if err != nil {
			return nil, fmt.Errorf("could not parse Q: %w", err)
		}
		fbEQ.Filters[i] = &FixedBandFilter{
			Frequency: freq,
			Gain:      gain,
			Q:         q,
		}
		i++
	}

	fbEQ.Filters = fbEQ.Filters[:i]

	return fbEQ, nil
}
