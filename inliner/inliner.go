package inliner

import (
	"code.google.com/p/go.net/html"
	"github.com/tbuckley/vulcanize/htmlutils"
	"github.com/tbuckley/vulcanize/pathresolver"
	"path/filepath"
	"regexp"
)

func IsExcluded(path string, excludes []*regexp.Regexp) bool {
	for _, pattern := range excludes {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

func InlineScripts(doc *htmlutils.Fragment, outputDir string, excludes []*regexp.Regexp) {
	// script:not([type])[src], script[type="text/javascript"][src]
	pred := htmlutils.AndP(
		htmlutils.HasTagnameP("script"),
		htmlutils.HasAttrP("src"),
		htmlutils.OrP(
			htmlutils.NotP(htmlutils.HasAttrP("type")),
			htmlutils.HasAttrValueP("type", "text/javascript")))

	scripts := doc.Search(pred)
	for _, script := range scripts {
		src, ok := htmlutils.Attr(script, "src")
		if ok && !isExcluded(src, excludes) {
			filename := filepath.Join(outputDir, src)
			content := ioutil.ReadFile(filename)
			inlinedScript := htmlutils.CreateScript(content)
			// @TODO: modify script content?
			htmlutils.ReplaceNodeWithNode(script, inlinedScript)
		}
	}
}

func InlineSheets(doc *htmlutils.Fragment, outputDir string, excludes []*regexp.Regexp) {
	// link[rel="stylesheet"]
	pred := htmlutils.AndP(htmlutils.HasTagnameP("link"), htmlutils.HasAttrValueP("rel", "stylesheet"))

	sheets := doc.Search(pred)
	for _, sheet := range sheets {
		href, ok := htmlutils.Attr(sheet, "href")
		if ok && !IsExcluded(href, excludes) {
			filename := filepath.Join(outputDir, src)
			content := ioutil.ReadFile(filename)
			content = pathresolver.RewriteURL(filepath.Dir(filename), outputDir, content)
			inlinedSheet := htmlutils.CreateStyle(content)
			// @TODO: copy link attributes (except rel/href) to style
			htmlutils.ReplaceNodeWithNode(sheet, inlinedSheet)
		}
	}
}
