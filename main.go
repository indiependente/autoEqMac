package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
	"github.com/indiependente/autoEqMac/server"
	au "github.com/logrusorgru/aurora"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	autoEQResults   = `https://raw.githubusercontent.com/jaakkopasanen/AutoEq/master/results`
	fixedBandSuffix = `%20FixedBandEQ.txt`
)

var (
	app  = kingpin.New("autoEqMac", "EqMac preset generator powered by AutoEq.\n\nAn interactive CLI that retrieves headphones EQ data from the AutoEq project and produces a JSON preset ready to be imported into EqMac.")
	file = app.Flag("file", "Output file path. By default it's the name of the headphones model selected.").Short('f').String()
	_    = app.HelpFlag.Short('h')
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
		return fmt.Errorf("⛔️ could not get EQ metadata: %w", err)
	}
	fmt.Println((au.Bold(au.Magenta("🎧 autoEqMac - EqMac preset generator powered by AutoEq"))))
	fmt.Println(au.Italic("Please select headphones model:"))
	headphones := prompt.Input("🎧 >>> ", populatedCompleter(eqMetas),
		prompt.OptionTitle("autoEqMac"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))
	fmt.Printf("👉 You selected: %s\n", headphones)

	eqMeta, err := srv.GetEQMetadataByName(headphones)
	if err != nil {
		return fmt.Errorf("⛔️ could not find EQ data for headphones %s: %w", headphones, err)
	}

	eqPreset, err := srv.GetFixedBandEQPreset(eqMeta.ID)
	if err != nil {
		return fmt.Errorf("⛔️ could not find fixed band EQ preset: %w", err)
	}

	filename := *file
	if filename == "" {
		filename = strings.ReplaceAll(headphones, " ", "_") + ".json"
	}
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("⛔️ could not create preset file: %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("⛔️ could not close preset file: %w", err)
	}

	err = srv.WritePreset(f, eqPreset)
	if err != nil {
		return fmt.Errorf("⛔️ could not write preset to file: %w", err)
	}
	fmt.Printf("📝 Preset saved to %s\n", f.Name())
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
