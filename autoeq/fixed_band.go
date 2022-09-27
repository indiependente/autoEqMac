package autoeq

import (
	"fmt"
	"strconv"
	"strings"
)

// FixedBandEQ represents a single EQ band.
type FixedBandEQ struct {
	Frequency int     // Hz
	Gain      float64 // Db
	Q         float64 // fixed
}

// FixedBandEQs represents a simple fixed bands EQ.
type FixedBandEQs []FixedBandEQ

// ToFixedBandEQs transforms the raw EQ data into a fixed band EQ.
// Returns an error if any.
func ToFixedBandEQs(data []byte) (FixedBandEQs, error) {
	var eqs FixedBandEQs
	rows := strings.Split(string(data), "\n")
	for _, row := range rows {
		if row == "" {
			continue
		}
		if strings.HasPrefix(row, "Preamp") {
			continue
		}
		eqFields := strings.Fields(row)
		if len(eqFields) < 12 {
			return nil, fmt.Errorf("could not parse : %s", row)
		}
		freq, err := strconv.Atoi(strings.TrimSpace(eqFields[5]))
		if err != nil {
			return nil, fmt.Errorf("could not parse frequency: %w", err)
		}
		gain, err := strconv.ParseFloat(strings.TrimSpace(eqFields[8]), 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse gain: %w", err)
		}
		q, err := strconv.ParseFloat(strings.TrimSpace(eqFields[11]), 64)
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
