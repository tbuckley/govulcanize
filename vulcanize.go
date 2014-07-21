package vulcanize

import (
	"bytes"
	"github.com/tbuckley/vulcanize/htmlutils"
	"github.com/tbuckley/vulcanize/pathresolver"
	"log"
	"os"
	"path"
	"path/filepath"

	"code.google.com/p/go.net/html"
)

type Options struct {
	Input     string
	Inline    bool
	OutputDir string
	Verbose   bool
}

var (
	options *Options
	read    map[string]bool = make(map[string]bool)
)

func HandleMainDocument() error {
	inputContents, err := os.Open(options.Input)
	if err != nil {
		return err
	}

	input, err := html.Parse(inputContents)
	if err != nil {
		return err
	}

	dir := path.Dir(options.Input)
	pathresolver.ResolvePaths(input, dir, options.OutputDir)

	processImports(input, true)
}

func processImports(input *html.Node, mainDoc bool) {
	matches := htmlutils.Search(input, htmlutils.IsImport)
	for _, el := range matches {
		href := htmlutils.Attr(el, "href")
		if !excludeImport(href) {
			importContent := concat(filepath.Join(options.OutputDir, href))
			if mainDoc {
				importContent = "<div hidden>" + importContent + "</div>"
			}
			ns, _ := htmlutils.ParseFragment(importContent)
			htmlutils.Replace(el, ns)
		}
	}
}

func excludeImport(href string) {

}

func concat(filename string) string {
	if !read[filename] {
		read[filename] = true
		doc := html.Parse(r)
		dir := path.Dir(filename)
		pathresolver.ResolvePaths(doc, dir, options.OutputDir)
		processImports(doc, false)
		buf := new(bytes.Buffer)
		html.Render(buf, doc)
		return buf.String()
	} else {
		if options.Verbose {
			log.Println("Dependency deduplicated")
		}
	}
}
