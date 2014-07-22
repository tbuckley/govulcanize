package importer

import (
	"github.com/tbuckley/vulcanize/htmlutils"
	"testing"
)

func TestNewImporter(t *testing.T) {

}

func TestImporter_Flatten(t *testing.T) {
	fragment, err := htmlutils.FromFile("../test/index.html")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(fragment.String())
	t.Error(fragment.String())
}

func TestImporter_load(t *testing.T) {

}

func TestImporter_processImports(t *testing.T) {

}

func TestImporter_excludeImport(t *testing.T) {

}

func TestImporter_deduplicateImport(t *testing.T) {

}
