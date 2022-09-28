// Package eqmac exports the equalization presets and maps autoEQ models.
package eqmac

// EQPreset represents an EqMac preset.
type EQPreset struct {
	Gains     Gains  `json:"gains"`
	ID        string `json:"id"`
	IsDefault bool   `json:"isDefault"`
	Name      string `json:"name"`
}

// Gains represents the dB gains of an EQPreset.
type Gains struct {
	Global float64   `json:"global"`
	Bands  []float64 `json:"bands"`
}
