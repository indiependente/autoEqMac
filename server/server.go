package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
)

const (
	headphonesIndex = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results/INDEX.md`
	// ErrEQMetadataNameNotFound returned when eq metadata name cannot be found.
	ErrEQMetadataNameNotFound Error = `eq metadata name not found`
	// ErrEQMetadataNotFound returned when eq metadata cannot be found.
	ErrEQMetadataNotFound Error = `eq metadata not found`
)

// Error is the error type returned by the server.
type Error string

// Error returns the string representation of an error.
func (e Error) Error() string {
	return string(e)
}

// HTTPServer is an HTTP implementation of a Server.
// It fulfills the requests received by obtaining data via HTTP requests.
type HTTPServer struct {
	client   Doer
	mdparser autoeq.MarkDownParser
	eqGetter autoeq.EQGetter
	mapper   mapping.Mapper
	eqMetas  map[string]*autoeq.EQMetadata
	eqNameID map[string]string
}

// NewHTTPServer returns a new HTTPServer.
func NewHTTPServer(d Doer, mdp autoeq.MarkDownParser, eqg autoeq.EQGetter, m mapping.Mapper) HTTPServer {
	return HTTPServer{
		client:   d,
		mdparser: mdp,
		eqGetter: eqg,
		mapper:   m,
		eqMetas:  map[string]*autoeq.EQMetadata{},
		eqNameID: map[string]string{},
	}
}

// ListEQsMetadata returns a list of all the EQ metadata found by the server.
// Returns an error if any.
func (s *HTTPServer) ListEQsMetadata() ([]*autoeq.EQMetadata, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, headphonesIndex, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get updated headphones list: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read headphones list raw data: %w", err)
	}
	eqMetas, err := s.mdparser.ParseMetadata(data)
	if err != nil {
		return nil, fmt.Errorf("could not parse headphones metadata: %w", err)
	}
	for _, meta := range eqMetas {
		s.eqMetas[meta.ID] = meta
		s.eqNameID[meta.Name] = meta.ID
	}

	return eqMetas, nil
}

// GetFixedBandEQPreset returns the EQ preset assiciated to the input id.
// Returns an error if any.
func (s *HTTPServer) GetFixedBandEQPreset(id string) (eqmac.EQPreset, error) {
	eqMeta, ok := s.eqMetas[id]
	if !ok {
		return eqmac.EQPreset{}, ErrEQMetadataNotFound
	}
	rawEQ, err := s.eqGetter.GetEQ(eqMeta)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not get raw EQ data: %w", err)
	}
	fbEQ, err := autoeq.ToFixedBandEQs(rawEQ)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not map raw EQ data: %w", err)
	}
	if fbEQ.Preamp == 0 { // fallback mechanism for Preamp
		globalPreamp, err := s.eqGetter.GetFixedBandGlobalPreamp(eqMeta) //nolint:govet // error shadowing doesn't impact its handling.
		if err != nil {
			return eqmac.EQPreset{}, fmt.Errorf("could not get global EQ preamp data: %w", err)
		}
		fbEQ.Preamp = globalPreamp
	}
	eqPreset, err := s.mapper.MapFixedBand(fbEQ, eqMeta)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not map raw EQ datato eqMac preset: %w", err)
	}

	return eqPreset, nil
}

// GetEQMetadataByName returns EQ metadata associated to a device name.
// Returns an error if any.
func (s *HTTPServer) GetEQMetadataByName(name string) (*autoeq.EQMetadata, error) {
	metaID, ok := s.eqNameID[name]
	if !ok {
		return nil, ErrEQMetadataNameNotFound
	}
	eqMeta, ok := s.eqMetas[metaID]
	if !ok {
		return nil, ErrEQMetadataNotFound
	}

	return eqMeta, nil
}

// WritePreset writes the input preset to the provided io.Writer.
// It uses json encoding.
// Returns an error if any.
func (s *HTTPServer) WritePreset(w io.Writer, p eqmac.EQPreset) error {
	jsonPreset, err := json.Marshal([]eqmac.EQPreset{p})
	if err != nil {
		return fmt.Errorf("could not marshal preset to JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(jsonPreset))

	return err
}
