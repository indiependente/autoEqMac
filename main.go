package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
)

const (
	headphonesIndex = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results/INDEX.md`
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	resp, err := http.Get(headphonesIndex)
	if err != nil {
		return fmt.Errorf("could not get updated headphones list: %w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read headphones list raw data: %w", err)
	}

	var parser autoeq.MarkDownParser = autoeq.MetadataParser{
		LinkPrefix:        "https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results",
		FixedBandEQSuffix: "%20FixedBandEQ.txt",
	}
	eqMetas, err := parser.ParseMetadata(data)
	if err != nil {
		return fmt.Errorf("could not parse headphones metadata: %w", err)
	}
	eqGetter := autoeq.EQHTTPGetter{
		Client: http.DefaultClient,
	}
	eqMeta := eqMetas[500]
	rawEQ, err := eqGetter.GetEQ(eqMeta)
	if err != nil {
		return fmt.Errorf("could not get raw EQ data: %w", err)
	}
	globalPreamp, err := eqGetter.GetFixedBandGlobalPreamp(eqMeta)
	if err != nil {
		return fmt.Errorf("could not get global EQ preamp data: %w", err)
	}
	eqMeta.Global = globalPreamp
	fbEQs, err := autoeq.ToFixedBandEQs(rawEQ)
	if err != nil {
		return fmt.Errorf("could not map raw EQ data: %w", err)
	}
	mapper := mapping.AutoEQMapper{}
	eqPreset, err := mapper.MapFixedBand(fbEQs, eqMeta)
	if err != nil {
		return fmt.Errorf("could not map raw EQ datato eqMac preset: %w", err)
	}
	jsonPreset, err := json.Marshal([]eqmac.EQPreset{eqPreset})
	if err != nil {
		return fmt.Errorf("could not write preset to JSON: %w", err)
	}
	fmt.Printf("%s\n", string(jsonPreset))
	return nil
}
