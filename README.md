# govulcanize

A golang version of the Vulcanize tool for Polymer.

### Dependencies

* code.google.com/p/go.net/html
* github.com/docopt/docopt.go

### Differences from nodejs version

* Whitespace
* Go's html package doesn't handle boolean attributes well (eg. `<div hidden="">` vs. `<div hidden>`)
* Go's html package will return a single node for the main document (so `<!doctype html>` may be lost)
* Go's html package doesn't use a trailing slash for singleton tags (eg. `<br>` instead of `<br/>`)

### TODO

* Enable `--strip` flag
* Handle deduplication of absolute/excluded imports
