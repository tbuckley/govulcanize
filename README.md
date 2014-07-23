# govulcanize

A golang version of the Vulcanize tool for Polymer.

Main differences in code output:

* Whitespace
* Go's html package doesn't handle boolean attributes well (eg. `<div hidden="">` vs. `<div hidden>`)
* Go's html package will return a single node for the main document (so `<!doctype html>` may be lost)
