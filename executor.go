package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/indiependente/autoEqMac/server"
)

// NewExecutor returns a new prompt.Executor.
func NewExecutor(srv *server.HTTPServer) func(string) {
	return func(s string) {
		s = strings.TrimSpace(s)
		if s == "bye" || s == "quit" {
			fmt.Println("Bye!")
			os.Exit(0)

			return
		}
		headphones := s

		fmt.Printf("👉 You selected: %s\n", headphones)
		eqMeta, err := srv.GetEQMetadataByName(headphones)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "⛔️ could not find EQ data for headphones %s: %q\n", headphones, err)
			os.Exit(1)

			return
		}

		eqPreset, err := srv.GetFixedBandEQPreset(eqMeta.ID)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "⛔️ could not find fixed band EQ preset: %q\n", err)
			os.Exit(1)

			return
		}

		filename := filepath.Clean(filename(file, headphones))
		f, err := os.Create(filename)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "⛔️ could not create preset file: %q\n", err)
			os.Exit(1)

			return
		}
		defer f.Close() //nolint: errcheck,gosec

		err = srv.WritePreset(f, eqPreset)
		if err != nil {
			_ = f.Close()
			_ = os.Remove(filename)
			_, _ = fmt.Fprintf(os.Stderr, "⛔️ could not write preset to file: %q\n", err)
			os.Exit(1)

			return
		}
		fmt.Printf("📝 Preset saved to %s\n", f.Name())

		os.Exit(0)
	}
}
