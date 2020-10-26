package autoeq

import (
	"fmt"
	"strconv"
	"strings"
)

type FixedBandEQ struct {
	Frequency int     // Hz
	Gain      float64 // Db
	Q         float64 // fixed
}

type FixedBandEQs []FixedBandEQ

func ToFixedBandEQs(data []byte) (FixedBandEQs, error) {
	var eqs FixedBandEQs
	rows := strings.Split(string(data), "\n")
	for _, row := range rows {
		eqFields := strings.Fields(row)
		freq, err := strconv.Atoi(eqFields[5])
		if err != nil {
			return nil, fmt.Errorf("could not parse frequency: %w", err)
		}
		gain, err := strconv.ParseFloat(eqFields[8], 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse gain: %w", err)
		}
		q, err := strconv.ParseFloat(eqFields[11], 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse Q: %w", err)
		}
		eqs = append(eqs, FixedBandEQ{
			Frequency: freq,
			Gain:      gain,
			Q:         q,
		})
	}
	return eqs, nil
}
