# govulcanize

A golang version of the Vulcanize tool for Polymer.

Main differences in code output:

* Go's html package doesn't handle boolean attributes well (eg. `&lt;div hidden=""&gt;` vs. `&lt;div hidden&gt;`)
* Go's html package will return a single node for the main document (so `&lt;!doctype html&gt;` may be lost)
