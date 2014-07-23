package optparser

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt.go"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var (
	DEFAULT_FILENAME = "vulcanized.html"
	ABS_URL          = regexp.MustCompilePOSIX("(^data:)|(^http[s]?:)|(^\\/)")
)

type Options struct {
	Input     string
	Output    string
	OutputDir string
	Excludes  Excludes

	CSP     bool
	CSPFile string
	Inline  bool
	Strip   bool

	Verbose bool
}

type Excludes struct {
	Imports []*regexp.Regexp
	Scripts []*regexp.Regexp
	Styles  []*regexp.Regexp
}

type Config struct {
	Excludes ConfigExcludes `json:"excludes"`
}

type ConfigExcludes struct {
	Imports []string `json:"imports"`
	Scripts []string `json:"scripts"`
	Styles  []string `json:"styles"`
}

func Parse() (*Options, error) {
	options := new(Options)
	config := new(Config)

	// Parse the command-line args
	arguments := parseArgs()

	// Initial configuration
	options.Excludes.Imports = []*regexp.Regexp{ABS_URL}
	options.Excludes.Scripts = []*regexp.Regexp{ABS_URL}
	options.Excludes.Styles = []*regexp.Regexp{ABS_URL}

	// Set initial options
	options.Input = arguments["<input>"].(string)
	options.Verbose = arguments["--verbose"].(bool)
	options.Strip = arguments["--strip"].(bool)
	options.Inline = arguments["--inline"].(bool)

	// Handle output
	outputFile, ok := arguments["--output"].(string)
	if ok {
		options.Output = outputFile
	} else {
		options.Output = filepath.Join(filepath.Dir(options.Input), DEFAULT_FILENAME)
	}
	options.OutputDir = filepath.Dir(options.Output)

	// Handle CSP
	options.CSP = arguments["--csp"].(bool)
	if options.CSP {
		dir, htmlFile := filepath.Split(options.Output)
		jsFile := htmlFile[:len(htmlFile)-len(".html")] + ".js"
		options.CSPFile = filepath.Join(dir, jsFile)
	}

	// Try to parse config file
	if arguments["--config"] != nil {
		configData, err := ioutil.ReadFile(arguments["--config"].(string))
		if err != nil {
			return nil, fmt.Errorf("Config file not found!")
		}
		err = json.Unmarshal([]byte(configData), &config)
		if err != nil {
			return nil, fmt.Errorf("Malformed config JSON!")
		}
	}

	// Read excludes from config file
	for _, restr := range config.Excludes.Imports {
		re, err := regexp.CompilePOSIX(restr)
		if err != nil {
			return nil, fmt.Errorf("Malformed import exclude config")
		}
		options.Excludes.Imports = append(options.Excludes.Imports, re)
	}
	for _, restr := range config.Excludes.Scripts {
		re, err := regexp.CompilePOSIX(restr)
		if err != nil {
			return nil, fmt.Errorf("Malformed import exclude config")
		}
		options.Excludes.Scripts = append(options.Excludes.Scripts, re)
	}
	for _, restr := range config.Excludes.Styles {
		re, err := regexp.CompilePOSIX(restr)
		if err != nil {
			return nil, fmt.Errorf("Malformed import exclude config")
		}
		options.Excludes.Styles = append(options.Excludes.Styles, re)
	}

	return options, nil
}

func parseArgs() map[string]interface{} {
	usage := `Go Vulcanize.

Usage:
  vulcanize [options] <input>

Options:
  -h, --help                  Show this screen.
  -v, --verbose               Verbose mode.
  -o <file>, --output <file>  Output file name (defaults to vulcanized.html).
  --config <file>             Read a given config file.
  --strip                     Remove comments and empty text nodes.
  --csp                       Extract inline scripts to a separate file (uses <output file name>.js).
  --inline                    The opposite of CSP mode, inline all assets (script and css) into the document.`

	arguments, _ := docopt.Parse(usage, nil, true, "Go Vulcanize 0.0.1", false)
	return arguments
}
