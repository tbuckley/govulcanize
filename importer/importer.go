package importer

import (
	"path/filepath"
	"regexp"

	"github.com/tbuckley/vulcanize/htmlutils"
)

type Importer struct {
	read             map[string]bool
	excludedPatterns []*regexp.Regexp
}

// NewImporter creates a new importer using the list of excluded patterns
func NewImporter(excludedPatterns []*regexp.Regexp) *Importer {
	return &Importer{
		read:             make(map[string]bool),
		excludedPatterns: excludedPatterns,
	}
}

// Flatten flattens out all of the imports from a document
func (i *Importer) Flatten(filename string) (*htmlutils.Fragment, error) {
	doc, err := i.load(filename)
	if err != nil {
		return nil, err
	}
	err = i.processImports(doc, filename)
	return doc, err
}

// load returns an HTML fragment representing the contents of the given file
// and ensures that the same file isn't loaded multiple times
func (i *Importer) load(filename string) (*htmlutils.Fragment, error) {
	doc, err := htmlutils.FromFile(filename)
	if err != nil {
		return nil, err
	}
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
				content, err := i.Flatten(importFile)
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
