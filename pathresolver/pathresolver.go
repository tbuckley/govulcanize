package pathresolver

import (
	"github.com/tbuckley/vulcanize/htmlutils"
	"path/filepath"
	"regexp"
)

var (
	ABS_URL      = regexp.MustCompilePOSIX("(^data:)|(^http[s]?:)|(^\\/)")
	URL          = regexp.MustCompilePOSIX("url\\([^)]*\\)")
	QUOTES       = regexp.MustCompilePOSIX("[\"']")
	URL_TEMPLATE = regexp.MustCompilePOSIX("{{.*}}")
)

func ResolvePaths(input *htmlutils.Fragment, inputPath string, outputPath string) {
	resolveAttributePaths(input, inputPath, outputPath)
	resolveCSSPaths(input, inputPath, outputPath)
	addAssetpathAttribute(input, inputPath, outputPath)
}

// resolveAttributePaths rewrites any relative URLs found in node attributes
// (eg. href, src, action, style)
func resolveAttributePaths(input *htmlutils.Fragment, inputPath string, outputPath string) {
	URL_ATTR := []string{"href", "src", "action", "style"}
	matches := input.Search(htmlutils.HasAnyAttrP(URL_ATTR...))
	for _, match := range matches {
		for _, attr := range URL_ATTR {
			if val, ok := htmlutils.Attr(match, attr); ok {
				if URL_TEMPLATE.FindAllStringIndex(val, 0) == nil {
					if attr == "style" {
						htmlutils.SetAttr(match, attr, rewriteURL(inputPath, outputPath, val))
					} else {
						htmlutils.SetAttr(match, attr, rewriteRelPath(inputPath, outputPath, val))
					}
				}
			}
		}
	}
}

// resolveCSSPaths rewrites any relative URLs found in CSS blocks
func resolveCSSPaths(input *htmlutils.Fragment, inputPath string, outputPath string) {
	matches := input.Search(htmlutils.IsStyleBlock)
	for _, match := range matches {
		text := rewriteURL(inputPath, outputPath, htmlutils.TextContent(match))
		htmlutils.SetTextContent(match, text)
	}
}

// addAssetpathAttribute adds the assetpath attribute to any polymer-element
// nodes that may be missing it
func addAssetpathAttribute(input *htmlutils.Fragment, inputPath string, outputPath string) {
	assetPath, _ := filepath.Rel(outputPath, inputPath)
	if assetPath != "" {
		assetPath += "/"
	}
	matches := input.Search(htmlutils.IsPolymerElementMissingAssetpath)
	for _, match := range matches {
		htmlutils.SetAttr(match, "assetpath", assetPath)
	}
}

// rewriteRelPath rewrites a path relative to inputPath to be relative to outputPath
func RewriteRelPath(inputPath string, outputPath string, rel string) string {
	if isAbsoluteURL(rel) {
		return rel
	}
	abs := filepath.Join(inputPath, rel)
	relPath, _ := filepath.Rel(outputPath, abs)
	return relPath
}

// rewriteURL converts all instances of `url('<RELPATH>')` in a CSS string to urls
// relative to the outputPath
func RewriteURL(inputPath string, outputPath string, cssText string) string {
	return URL.ReplaceAllStringFunc(cssText, func(match string) string {
		path := stripQuotes(match)
		path = path[4 : len(path)-1]
		path = rewriteRelPath(inputPath, outputPath, path)
		return "url(" + path + ")"
	})
}

// isAbsoluteURL returns true if url is absolute
func isAbsoluteURL(url string) bool {
	return ABS_URL.MatchString(url)
}

// stripQuotes removes all single and double quotes from a string
func stripQuotes(str string) string {
	return QUOTES.ReplaceAllString(str, "")
}
