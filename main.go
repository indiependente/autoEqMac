package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/google/uuid"
	au "github.com/logrusorgru/aurora"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/indiependente/autoEqMac/autoeq"
	"github.com/indiependente/autoEqMac/eqmac/mapping"
	"github.com/indiependente/autoEqMac/server"
)

var (
	app  = kingpin.New("autoEqMac", "EqMac preset generator powered by AutoEq.\n\nAn interactive CLI that retrieves headphones EQ data from the AutoEq project and produces a JSON preset ready to be imported into EqMac.")
	file = app.Flag("file", "Output file path. By default it's the name of the headphones model selected.").Short('f').String()
	_    = app.Version(fmt.Sprintf("autoEqMac %s commit %s built by %s on %s", version, commit, builtBy, date)).VersionFlag.Short('v')
	_    = app.HelpFlag.Short('h')

	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
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
	mdParser := autoeq.NewMetadataParser()
	eqGetter := autoeq.EQHTTPGetter{
		Client: http.DefaultClient,
	}
	mapper := mapping.NewAutoEQMapper(mapping.WrappedGenerator(func() string {
		return uuid.New().String()
	}))

	srv := server.NewHTTPServer(client, mdParser, eqGetter, mapper)
	eqMetas, err := srv.ListEQsMetadata()
	if err != nil {
		return fmt.Errorf("⛔️ could not get EQ metadata: %w", err)
	}
	fmt.Println(au.Bold(au.Magenta("🎧 autoEqMac - EqMac preset generator powered by AutoEq")))
	fmt.Println(au.Italic("Please select headphones model:"))

	p := prompt.New(NewExecutor(&srv),
		populatedCompleter(eqMetas),
		prompt.OptionTitle("autoEqMac"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator))
	p.Run()

	return nil
}

func filename(file *string, headphones string) string {
	filename := *file
	if filename == "" {
		filename = fmt.Sprintf("%s.json", strings.ReplaceAll(headphones, " ", "_"))
	}
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	return filename
}
