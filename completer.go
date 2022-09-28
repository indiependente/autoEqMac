package main

import (
	"github.com/c-bata/go-prompt"

	"github.com/indiependente/autoEqMac/autoeq"
)

func populatedCompleter(eqMetas []*autoeq.EQMetadata) func(prompt.Document) []prompt.Suggest {
	suggs := make([]prompt.Suggest, len(eqMetas))
	for i, meta := range eqMetas {
		suggs[i] = prompt.Suggest{
			Text: meta.Name, Description: meta.Author,
		}
	}

	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggs, d.CurrentLine(), true)
	}
}
