package importer

import (
	"code.google.com/p/go.net/html"
	"github.com/tbuckley/vulcanize/htmlutils"
	"regexp"
	"testing"
)

func TestNewImporter(t *testing.T) {
	re1 := regexp.MustCompilePOSIX("href.*")
	re2 := regexp.MustCompilePOSIX("data.*")
	i := New([]*regexp.Regexp{re1, re2}, "./")

	if i == nil {
		t.Error("returned importer is null")
	}

	if len(i.excludedPatterns) != 2 {
		t.Error("returned importer does not have excluded patterns")
	}

	if i.read == nil {
		t.Error("returned importer does not have a read map")
	}
}

func TestImporter_Flatten(t *testing.T) {
	i := New(nil, "./")

	doc, err := i.Flatten("../test/index.html", nil)
	t.Log(doc.String())
	if err != nil {
		t.Error(err.Error())
	}

	var els []*html.Node

	els = doc.Search(htmlutils.HasTagnameP("foo-a"))
	if len(els) != 1 {
		t.Error("foo-a tag missing from vulcanized document")
	}

	els = doc.Search(htmlutils.AndP(htmlutils.HasTagnameP("polymer-element"), htmlutils.HasAttrValueP("name", "foo-a")))
	if len(els) != 1 {
		t.Error("polymer-element[name=\"foo-b\"] tag missing from vulcanized document")
	}

	els = doc.Search(htmlutils.AndP(htmlutils.HasTagnameP("polymer-element"), htmlutils.HasAttrValueP("name", "foo-b")))
	if len(els) != 1 {
		t.Error("polymer-element[name=\"foo-b\"] tag missing from vulcanized document")
	}

	els = doc.Search(htmlutils.HasTagnameP("link"))
	if len(els) != 0 {
		t.Error("link tags left in vulcanized document")
	}
}

func TestImporter_load(t *testing.T) {

}

func TestImporter_processImports(t *testing.T) {

}

func TestImporter_excludeImport(t *testing.T) {

}

func TestImporter_deduplicateImport(t *testing.T) {

}
