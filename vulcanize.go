package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tbuckley/vulcanize/htmlutils"
	"github.com/tbuckley/vulcanize/importer"
	"github.com/tbuckley/vulcanize/inliner"
	"github.com/tbuckley/vulcanize/optparser"
)

func main() {
	var err error

	// Parse options
	options, err := optparser.Parse()
	handleError(err)

	// Import doc
	importer := importer.New(options.Excludes.Imports, options.Excludes.Styles, options.OutputDir)
	doc, err := importer.Flatten(options.Input, nil)
	handleError(err)

	// Messy logic for inlining and handling csp
	if options.Inline {
		err := inliner.InlineScripts(doc, options.OutputDir, options.Excludes.Scripts)
		handleError(err)
	}
	UseNamedPolymerInvocations(doc, options.Verbose)
	if options.CSP {
		SeparateScripts(doc, options.CSPFile, options.Verbose)
	}

	// Clean up
	DeduplicateImports(doc)
	if options.Strip {
		RemoveCommentsAndWhitespace(doc)
	}

	WriteFile(doc, options.Output)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(-1)
	}
}

func UseNamedPolymerInvocations(doc *htmlutils.Fragment, verbose bool) {
	// script:not([type]):not([src]), script[type="text/javascript"]:not([src])
	pred := htmlutils.AndP(
		htmlutils.HasTagnameP("script"),
		htmlutils.NotP(htmlutils.HasAttrP("src")),
		htmlutils.OrP(
			htmlutils.NotP(htmlutils.HasAttrP("type")),
			htmlutils.HasAttrValueP("type", "text/javascript")))

	POLYMER_INVOCATION := regexp.MustCompile("Polymer\\(([^,{]+)?(?:,\\s*)?({|\\))")
	inlineScripts := doc.Search(pred)
	for _, script := range inlineScripts {
		content := htmlutils.TextContent(script)
		parentElement := htmlutils.Closest(script, htmlutils.HasTagnameP("polymer-element"))
		if parentElement != nil {
			match := POLYMER_INVOCATION.FindStringSubmatch(content)
			if len(match) != 0 && match[1] == "" {
				name, _ := htmlutils.Attr(parentElement, "name")
				// @TODO handle case where name is not defined
				namedInvocation := "Polymer('" + name + "'"
				if match[2] == "{" {
					namedInvocation += ",{"
				} else {
					namedInvocation += ")"
				}
				content = strings.Replace(content, match[0], namedInvocation, 1)
				if verbose {
					fmt.Printf("%s -> %s\n", match[0], namedInvocation)
				}
				htmlutils.SetTextContent(script, content)
			}
		}
	}
}

func SeparateScripts(doc *htmlutils.Fragment, filename string, verbose bool) {
	if verbose {
		fmt.Println("Separating scripts into separate file")
	}

	// script:not([type]):not([src]), script[type="text/javascript"]:not([src])
	pred := htmlutils.AndP(
		htmlutils.HasTagnameP("script"),
		htmlutils.NotP(htmlutils.HasAttrP("src")),
		htmlutils.OrP(
			htmlutils.NotP(htmlutils.HasAttrP("type")),
			htmlutils.HasAttrValueP("type", "text/javascript")))

	inlineScripts := doc.Search(pred)
	scripts := make([]string, 0, len(inlineScripts))
	for _, script := range inlineScripts {
		content := htmlutils.TextContent(script)
		scripts = append(scripts, content)
		htmlutils.RemoveNode(doc, script)
	}

	scriptContent := strings.Join(scripts, ";\n")
	// @TODO compress if --strip is set
	ioutil.WriteFile(filename, []byte(scriptContent), 0775)

	// insert out-of-lined script into document
	basename := filepath.Base(filename)
	script := htmlutils.CreateExternalScript(basename)
	matches := doc.Search(htmlutils.HasTagnameP("body"))
	// TODO ensure that len(matches) > 0
	body := matches[0]
	body.AppendChild(script)
}

func DeduplicateImports(doc *htmlutils.Fragment) {

}

func RemoveCommentsAndWhitespace(doc *htmlutils.Fragment) {

}

func WriteFile(doc *htmlutils.Fragment, filename string) {
	content := doc.String()
	ioutil.WriteFile(filename, []byte(content), 0775)
}
