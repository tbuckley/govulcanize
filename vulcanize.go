package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tbuckley/vulcanize/importer"
	"github.com/tbuckley/vulcanize/optparser"
)

func main() {
	var err error

	// Parse options
	options, err := optparser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Import doc
	importer := importer.New(options.Excludes.Imports, options.OutputDir)
	doc, err := importer.Flatten(options.Input, nil)
	if err != nil {
		return err
	}

	// Messy logic...
	if options.Inline {
		InlineScripts(doc, options.OutputDir)
	}
	UseNamedPolymerInvocations(doc)
	if options.CSP {
		SeparateScripts(doc)
	}
	DeduplicateImports(doc)
	if options.Strip {
		RemoveCommentsAndWhitespace(doc)
	}
	WriteFile(doc)

	HandleMainDocument()
}
