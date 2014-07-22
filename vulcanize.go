package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tbuckley/vulcanize/importer"
	"github.com/tbuckley/vulcanize/optparser"
)

var (
	options *optparser.Options
	logger  *log.Logger
)

func HandleMainDocument() error {
	i := importer.New(options.Excludes.Imports, options.OutputDir)
	doc, err := i.Flatten(options.Input, nil)
	if err != nil {
		return err
	}

	fmt.Println(doc.String())

	return nil
}

func main() {
	logger = log.New(os.Stdout, "logger:", log.Lshortfile)

	options, err := optparser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", options)
	// HandleMainDocument()
}
