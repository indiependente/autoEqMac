# autoEqMac
An interactive CLI that retrieves headphones EQ data from the [AutoEq Project](https://github.com/jaakkopasanen/AutoEq) and produces a JSON preset ready to be imported into [EqMac](https://github.com/bitgapp/eqMac/).

## How to

```
▶ ./autoEqMac --help
usage: autoEqMac [<flags>]

An interactive CLI that retrieves headphones EQ data from the AutoEq project and
produces a JSON preset ready to be imported into EqMac.

Flags:
      --help       Show context-sensitive help (also try --help-long and
                   --help-man).
  -f, --file=FILE  Output file path.
```

### Example usage

[![asciicast](https://asciinema.org/a/368415.svg)](https://asciinema.org/a/368415)

Once the JSON content has been generated and saved into a file, you can import it into eqMac.
