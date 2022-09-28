package autoeq

import (
	"fmt"
	"strconv"
	"strings"
)

const preampFields = 12

// FixedBandEQ represents a single EQ band.
type FixedBandEQ struct {
	Frequency int     // Hz
	Gain      float64 // Db
	Q         float64 // fixed
}

// FixedBandEQs represents a simple fixed bands EQ.
type FixedBandEQs []*FixedBandEQ

// ToFixedBandEQs transforms the raw EQ data into a fixed band EQ.
// Returns an error if any.
func ToFixedBandEQs(data []byte) (FixedBandEQs, error) {
	rows := strings.Split(string(data), "\n")
	eqs := make(FixedBandEQs, len(rows))

	i := 0
	for _, row := range rows {
		if row == "" {
			continue
		}
		if strings.HasPrefix(row, "Preamp") {
			continue
		}
		eqFields := strings.Fields(row)
		if len(eqFields) < preampFields {
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
		eqs[i] = &FixedBandEQ{
			Frequency: freq,
			Gain:      gain,
			Q:         q,
		}
		i++
	}

	return eqs[:i], nil
}
