package vulcanize

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"code.google.com/p/go.net/html"

	"github.com/tbuckley/vulcanize/htmlutils"
	"github.com/tbuckley/vulcanize/importer"
	"github.com/tbuckley/vulcanize/pathresolver"
)

type Options struct {
	Input     string   `json:"input"`
	Inline    bool     `json:"inline"`
	OutputDir string   `json:"output"`
	Verbose   bool     `json:"verbose"`
	Excludes  Excludes `json:"excludes"`
}

type Excludes struct {
	Imports []*regexp.Regexp `json:"imports"`
	Scripts []*regexp.Regexp `json:"scripts"`
	Styles  []*regexp.Regexp `json:"styles"`
}

var (
	options *Options
)

func HandleMainDocument() error {
	i := importer.NewImporter(options.Excludes.Imports)
	doc, err := i.Flatten(options.Input)
	if err != nil {
		return err
	}

	dir := path.Dir(options.Input)
	pathresolver.ResolvePaths(doc, dir, options.OutputDir)
}
