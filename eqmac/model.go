package eqmac

type EQPreset struct {
	Gains     Gains  `json:"gains"`
	ID        string `json:"id"`
	IsDefault bool   `json:"isDefault"`
	Name      string `json:"name"`
}
type Gains struct {
	Global float64   `json:"global"`
	Bands  []float64 `json:"bands"`
}
