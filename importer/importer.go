package importer

import (
	"code.google.com/p/go.net/html"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tbuckley/vulcanize/htmlutils"
	"github.com/tbuckley/vulcanize/pathresolver"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "logger:", log.Lshortfile)
}

type Importer struct {
	read             map[string]bool
	excludedPatterns []*regexp.Regexp
	outputDir        string
}

// NewImporter creates a new importer using the list of excluded patterns
func New(excludedPatterns []*regexp.Regexp, outputDir string) *Importer {
	return &Importer{
		read:             make(map[string]bool),
		excludedPatterns: excludedPatterns,
		outputDir:        outputDir,
	}
}

// Flatten flattens out all of the imports from a document
func (i *Importer) Flatten(filename string, context *html.Node) (*htmlutils.Fragment, error) {
	logger.Printf("Flatten: %v", filename)
	doc, err := i.load(filename, context)
	if err != nil {
		return nil, err
	}
	err = i.processImports(doc, filename)
	return doc, err
}

// load returns an HTML fragment representing the contents of the given file
// and ensures that the same file isn't loaded multiple times
func (i *Importer) load(filename string, context *html.Node) (*htmlutils.Fragment, error) {
	doc, err := htmlutils.FromFile(filename, context)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(filename)
	pathresolver.ResolvePaths(doc, dir, i.outputDir)

	i.read[filename] = true
	return doc, nil
}

// processImports iterates over the imports in a document, inlining available
// ones and skipping those that have been excluded
func (i *Importer) processImports(doc *htmlutils.Fragment, filename string) error {
	imports := doc.Search(htmlutils.IsImport)
	for _, imp := range imports {
		href, ok := htmlutils.Attr(imp, "href")
		if ok && !i.excludeImport(href) {
			dir := filepath.Dir(filename)
			importFile := filepath.Join(dir, href)
			if i.deduplicateImport(importFile) {
				htmlutils.RemoveNode(doc, imp)
			} else {
				content, err := i.Flatten(importFile, imp.Parent)
				if err != nil {
					return err
				}
				htmlutils.ReplaceNodeWithFragment(doc, imp, content)
			}
		}
	}
	return nil
}

// excludeImport returns true if the provided href should not be imported
func (i *Importer) excludeImport(href string) bool {
	for _, pattern := range i.excludedPatterns {
		if pattern.MatchString(href) {
			return true
		}
	}
	return false
}

// deduplicateImport returns true if filename has already been imported
func (i *Importer) deduplicateImport(filename string) bool {
	_, ok := i.read[filename]
	return ok
}
