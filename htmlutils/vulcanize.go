package htmlutils

import (
	"code.google.com/p/go.net/html"
)

func IsPolymerElementMissingAssetpath(n *html.Node) bool {
	_, hasAssetpath := Attr(n, "assetpath")
	return n.Type == html.ElementNode && n.Data == "polymer-element" && !hasAssetpath
}

// IsStyleBlock returns true if the given html node is a <style> node
// with type="text/css" or no type set
func IsStyleBlock(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "style" {
		kind, ok := Attr(n, "type")
		return !ok || kind == "text/css"
	}
	return false
}

// IsImport returns true if the given html node matches link[rel="import"][href]
func IsImport(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "link" {
		relType, relOk := Attr(n, "rel")
		_, hasHref := Attr(n, "href")
		return relOk && relType == "import" && hasHref
	}
	return false
}
