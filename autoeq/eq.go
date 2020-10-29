package autoeq

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	fixedBandEQTitle = "### Fixed Band EQs"
	preampPrefix     = "apply preamp of **"
)

// EQGetter defines the behaviour of a component capable of retrieving EQ information.
type EQGetter interface {
	GetEQ(meta EQMetadata) ([]byte, error)
	GetFixedBandGlobalPreamp(meta EQMetadata) (float64, error)
}

// Doer defines the behaviour of a component capable of doing an HTTP request,
// returning an HTTP response and an error.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// compile time interface implementation check
var _ EQGetter = EQHTTPGetter{}

// EQHTTPGetter is an HTTP based implementation of an EQGetter.
type EQHTTPGetter struct {
	Client Doer
}

// GetEQ returns the raw bytes that represent the EQ described by the input EQ metadata.
// Returns an error if any.
func (g EQHTTPGetter) GetEQ(meta EQMetadata) ([]byte, error) {
	rawdata, err := do(g.Client, meta.Link)
	if err != nil {
		return nil, fmt.Errorf("could not get eq: %w", err)
	}
	return rawdata, nil
}

// GetFixedBandGlobalPreamp returns the global preamp value in dB for the input EQ.
// Returns an error if any.
func (g EQHTTPGetter) GetFixedBandGlobalPreamp(meta EQMetadata) (float64, error) {
	rawdata, err := do(g.Client, meta.Link[:strings.LastIndex(meta.Link, "/")]+"/README.md")
	if err != nil {
		return 0, fmt.Errorf("could not get eq: %w", err)
	}
	//remove db
	globalRaw := strings.Replace(string(rawdata), "dB", "", -1)
	global, err := extractGlobalPreamp(globalRaw)
	if err != nil {
		return 0, fmt.Errorf("could not extract global preamp value: %w", err)
	}
	return global, nil
}

func extractGlobalPreamp(data string) (float64, error) {
	fbraw := strings.Split(data, fixedBandEQTitle)[1]
	dbraw := strings.Split(fbraw, preampPrefix)[1]
	return strconv.ParseFloat(strings.Split(dbraw, "**")[0], 64)
}

func do(doer Doer, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}
	resp, err := doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform HTTP request: %w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response data: %w", err)
	}
	return data, nil
}
