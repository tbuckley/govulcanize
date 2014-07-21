package pathresolver

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"github.com/tbuckley/vulcanize/htmlutils"
	"strings"
	"testing"
)

func TestPathResolver_resolveAttributePaths(t *testing.T) {
	p := new(PathResolver)

	helper := func(input string, inputPath string, outputPath string, id string) string {
		r := strings.NewReader(input)
		document, _ := html.Parse(r)
		p.resolveAttributePaths(document, inputPath, outputPath)
		buf := new(bytes.Buffer)
		target := htmlutils.GetElementById(document, id)
		if target != nil {
			html.Render(buf, target)
			return buf.String()
		}
		return ""
	}

	output := helper("<a id=\"target\" href=\"qux/page.html\"></a>", "/foo/bar", "/foo/baz", "target")
	expected := "<a id=\"target\" href=\"../bar/qux/page.html\"></a>"
	if output != expected {
		t.Errorf("Expected %v, got %v", expected, output)
	}

	output = helper("<form id=\"target\" action=\"qux/page.html\"></form>", "/foo/bar", "/foo/baz", "target")
	expected = "<form id=\"target\" action=\"../bar/qux/page.html\"></form>"
	if output != expected {
		t.Errorf("Expected %v, got %v", expected, output)
	}

	output = helper("<a id=\"target\" style=\"background-image: url('qux/page.html');\"></a>", "/foo/bar", "/foo/baz", "target")
	expected = "<a id=\"target\" style=\"background-image: url(../bar/qux/page.html);\"></a>"
	if output != expected {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestPathResolver_resolveCSSPaths(t *testing.T) {
	p := new(PathResolver)

	helper := func(input string, inputPath string, outputPath string, id string) string {
		r := strings.NewReader(input)
		document, _ := html.Parse(r)
		p.resolveCSSPaths(document, inputPath, outputPath)
		buf := new(bytes.Buffer)
		target := htmlutils.GetElementById(document, id)
		if target != nil {
			html.Render(buf, target)
			return buf.String()
		}
		return ""
	}

	output := helper("<style id=\"target\">body {background-image: url('qux/page.html');}</style>", "/foo/bar", "/foo/baz", "target")
	expected := "<style id=\"target\">body {background-image: url(../bar/qux/page.html);}</style>"
	if output != expected {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestPathResolver_addAssetpathAttribute(t *testing.T) {
	p := new(PathResolver)

	helper := func(input string, inputPath string, outputPath string, id string) string {
		r := strings.NewReader(input)
		document, _ := html.Parse(r)
		p.addAssetpathAttribute(document, inputPath, outputPath)
		buf := new(bytes.Buffer)
		target := htmlutils.GetElementById(document, id)
		if target != nil {
			html.Render(buf, target)
			return buf.String()
		}
		return ""
	}

	output := helper("<polymer-element id=\"target\"></polymer-element>", "/foo/bar", "/foo/baz", "target")
	expected := "<polymer-element id=\"target\" assetpath=\"../bar/\"></polymer-element>"
	if output != expected {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestPathResolver_rewriteRelPath(t *testing.T) {
	p := new(PathResolver)

	result := p.rewriteRelPath("/foo/bar", "/foo/baz", "qux/page.html")
	if result != "../bar/qux/page.html" {
		t.Errorf("Expected %v, got %v", "../bar/qux/page.html", result)
	}
}

func TestPathResolver_rewriteURL(t *testing.T) {
	p := new(PathResolver)
	cssText := p.rewriteURL("/foo/bar", "/foo/baz", "background-image: url('backgrounds/bkg.png')")
	if cssText != "background-image: url(../bar/backgrounds/bkg.png)" {
		t.Errorf("Expected rewritten url, got %v", cssText)
	}
}
