package inliner

import (
	"code.google.com/p/go.net/html"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/tbuckley/vulcanize/htmlutils"
	"github.com/tbuckley/vulcanize/pathresolver"
)

func IsExcluded(path string, excludes []*regexp.Regexp) bool {
	for _, pattern := range excludes {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

func InlineScripts(doc *htmlutils.Fragment, outputDir string, excludes []*regexp.Regexp) error {
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
		if ok && !IsExcluded(src, excludes) {
			filename := filepath.Join(outputDir, src)
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			inlinedScript := htmlutils.CreateScript(string(content))
			// @TODO: modify script content?
			htmlutils.ReplaceNodeWithNode(doc, script, inlinedScript)
		}
	}
	return nil
}

func InlineSheets(doc *htmlutils.Fragment, outputDir string, excludes []*regexp.Regexp) error {
	// link[rel="stylesheet"]
	pred := htmlutils.AndP(htmlutils.HasTagnameP("link"), htmlutils.HasAttrValueP("rel", "stylesheet"))

	sheets := doc.Search(pred)
	for _, sheet := range sheets {
		href, ok := htmlutils.Attr(sheet, "href")
		if ok && !IsExcluded(href, excludes) {
			filename := filepath.Join(outputDir, href)
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			stylesheet := string(content)
			stylesheet = pathresolver.RewriteURL(filepath.Dir(filename), outputDir, stylesheet)
			inlinedSheet := htmlutils.CreateStyle(stylesheet)
			// @TODO: copy link attributes (except rel/href) to style
			for _, attr := range sheet.Attr {
				if attr.Key != "rel" && attr.Key != "href" {
					inlinedSheet.Attr = append(inlinedSheet.Attr, html.Attribute{
						Key: attr.Key,
						Val: attr.Val,
					})
				}
			}
			htmlutils.ReplaceNodeWithNode(doc, sheet, inlinedSheet)
		}
	}
	return nil
}
