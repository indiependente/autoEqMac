package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
)

const (
	headphonesIndex                       = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results/INDEX.md`
	ErrEQMetadataNameNotFound ServerError = `eq metadata name not found`
	ErrEQMetadataNotFound     ServerError = `eq metadata not found`
)

type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPServer struct {
	client   Doer
	mdparser autoeq.MarkDownParser
	eqGetter autoeq.EQGetter
	mapper   mapping.Mapper
	eqMetas  map[string]autoeq.EQMetadata
	eqNameID map[string]string
}

func NewHTTPServer(d Doer, mdp autoeq.MarkDownParser, eqg autoeq.EQGetter, m mapping.Mapper) HTTPServer {
	return HTTPServer{
		client:   d,
		mdparser: mdp,
		eqGetter: eqg,
		mapper:   m,
		eqMetas:  map[string]autoeq.EQMetadata{},
		eqNameID: map[string]string{},
	}
}

func (s HTTPServer) ListEQsMetadata() ([]autoeq.EQMetadata, error) {
	resp, err := http.Get(headphonesIndex)
	if err != nil {
		return nil, fmt.Errorf("could not get updated headphones list: %w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
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

func (s HTTPServer) GetFixedBandEQPreset(id string) (eqmac.EQPreset, error) {
	eqMeta, ok := s.eqMetas[id]
	if !ok {
		return eqmac.EQPreset{}, ErrEQMetadataNotFound
	}
	rawEQ, err := s.eqGetter.GetEQ(eqMeta)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not get raw EQ data: %w", err)
	}
	globalPreamp, err := s.eqGetter.GetFixedBandGlobalPreamp(eqMeta)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not get global EQ preamp data: %w", err)
	}
	eqMeta.Global = globalPreamp
	fbEQs, err := autoeq.ToFixedBandEQs(rawEQ)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not map raw EQ data: %w", err)
	}
	eqPreset, err := s.mapper.MapFixedBand(fbEQs, eqMeta)
	if err != nil {
		return eqmac.EQPreset{}, fmt.Errorf("could not map raw EQ datato eqMac preset: %w", err)
	}
	return eqPreset, nil
}

func (s HTTPServer) GetEQMetadataByName(name string) (autoeq.EQMetadata, error) {
	metaID, ok := s.eqNameID[name]
	if !ok {
		return autoeq.EQMetadata{}, ErrEQMetadataNameNotFound
	}
	eqMeta, ok := s.eqMetas[metaID]
	if !ok {
		return autoeq.EQMetadata{}, ErrEQMetadataNotFound
	}
	return eqMeta, nil
}

func (s HTTPServer) WritePreset(w io.Writer, p eqmac.EQPreset) error {
	jsonPreset, err := json.Marshal([]eqmac.EQPreset{p})
	if err != nil {
		return fmt.Errorf("could not marshal preset to JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(jsonPreset))
	return err
}
