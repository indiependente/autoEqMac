//go:generate mockgen -package server -source=interface.go -destination server_mock.go

package server

import (
	"io"
	"net/http"

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

// Doer defines the behaviour of a component capable of doing an HTTP request,
// returning an HTTP response and an error.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
