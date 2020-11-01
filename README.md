[![Go Report Card](https://goreportcard.com/badge/github.com/indiependente/autoEqMac)](https://goreportcard.com/report/github.com/indiependente/autoEqMac)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-72%25-brightgreen.svg?longCache=true&style=flat)</a>
[![Workflow Status](https://github.com/indiependente/autoEqMac/workflows/lint-test/badge.svg)](https://github.com/indiependente/autoEqMac/actions)
# autoEqMac
An interactive CLI that retrieves headphones EQ data from the [AutoEq Project](https://github.com/jaakkopasanen/AutoEq) and produces a JSON preset ready to be imported into [EqMac](https://github.com/bitgapp/eqMac/).

## Dependencies
 - Go

## How to

### Install

`go get github.com/indiependente/autoEqMac`

### Supported commands

```
▶ autoEqMac --help
usage: autoEqMac [<flags>]

EqMac preset generator powered by AutoEq.

An interactive CLI that retrieves headphones EQ data from the AutoEq project and 
produces a JSON preset ready to be imported into EqMac.

Flags:
      --help       Show context-sensitive help (also try --help-long and
                   --help-man).
  -f, --file=FILE  Output file path. By default it's the name of the headphones
                   model selected.
```

### Example usage

[![asciicast](https://asciinema.org/a/368884.svg)](https://asciinema.org/a/368884)

Once the JSON content has been generated and saved into a file, you can import it into eqMac.

By default `autoEqMac` saves a JSON file with the same name of the headphones model you selected in the current working directory.

You can provide a different path by passing it using the `-f, --file` flag.

## TODO
- [ ] GUI

## Credits

Thanks to:
 - https://github.com/jaakkopasanen/AutoEq
 - https://github.com/bitgapp/eqMac/
 - https://github.com/c-bata/go-prompt
