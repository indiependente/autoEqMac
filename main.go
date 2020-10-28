package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
	"github.com/indiependente/autoEqMac/server"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	autoEQResults   = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results`
	fixedBandSuffix = `%20FixedBandEQ.txt`
)

var (
	app  = kingpin.New("autoEqMac", "An interactive CLI that retrieves headphones EQ data from the AutoEq project and produces a JSON preset ready to be imported into EqMac.")
	file = app.Flag("file", "Output file path.").Short('f').String()
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	client := http.DefaultClient
	mdParser := autoeq.MetadataParser{
		LinkPrefix:        autoEQResults,
		FixedBandEQSuffix: fixedBandSuffix,
	}
	eqGetter := autoeq.EQHTTPGetter{
		Client: http.DefaultClient,
	}
	mapper := mapping.AutoEQMapper{}
	srv := server.NewHTTPServer(client, mdParser, eqGetter, mapper)
	eqMetas, err := srv.ListEQsMetadata()
	if err != nil {
		return fmt.Errorf("could not get EQ metadata: %w", err)
	}

	fmt.Println("Please select headphones model:")
	t := prompt.Input("🎧 >>> ", populatedCompleter(eqMetas),
		prompt.OptionTitle("autoEqMac"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))
	fmt.Println("You selected " + t)

	eqMeta, err := srv.GetEQMetadataByName(t)
	if err != nil {
		return fmt.Errorf("could not find EQ data for headphones %s: %w", t, err)
	}

	eqPreset, err := srv.GetFixedBandEQPreset(eqMeta.ID)
	if err != nil {
		return fmt.Errorf("could not find fixed band EQ preset: %w", err)
	}
	out := os.Stdout
	if *file != "" {
		f, err := os.Create(*file)
		if err != nil {
			return fmt.Errorf("could not create preset file: %w", err)
		}
		out = f
	}
	err = srv.WritePreset(out, eqPreset)
	if err != nil {
		return fmt.Errorf("could not write preset to file: %w", err)
	}
	return nil
}

func populatedCompleter(eqMetas []autoeq.EQMetadata) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		var suggs []prompt.Suggest
		for _, meta := range eqMetas {
			suggs = append(suggs, prompt.Suggest{
				Text: meta.Name, Description: meta.ID,
			})
		}
		return prompt.FilterContains(suggs, d.Text, true)
	}
}
