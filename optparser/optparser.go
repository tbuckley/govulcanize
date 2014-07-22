package optparser

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var (
	input      = flag.String("input", "", "Input file name")
	output     = flag.String("output", "vulcanized.html", "Output file name")
	verbose    = flag.Bool("verbose", false, "More verbose logging")
	help       = flag.Bool("help", false, "Print this message")
	configFile = flag.String("config", "", "Read a given config file")
	strip      = flag.Bool("strip", false, "Remove comments and empty text nodes")
	csp        = flag.Bool("csp", false, "Extract inline scripts to a separate file (uses <output file name>.js)")
	inline     = flag.Bool("inline", false, "The opposite of CSP mode, inline all assets (script and css) into the document")

	DEFAULT_FILENAME = "vulcanized.html"

	ABS_URL = regexp.MustCompilePOSIX("(^data:)|(^http[s]?:)|(^\\/)")
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
	Help    bool
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
	flag.Parse()

	// Input file may also be a positional arg
	if flag.NArg() == 1 {
		*input = flag.Arg(0)
	}

	options := new(Options)
	config := new(Config)

	options.Excludes.Imports = []*regexp.Regexp{ABS_URL}
	options.Excludes.Scripts = []*regexp.Regexp{ABS_URL}
	options.Excludes.Styles = []*regexp.Regexp{ABS_URL}

	if *configFile != "" {
		configData, err := ioutil.ReadFile(*configFile)
		if err != nil {
			return nil, fmt.Errorf("Config file not found!")
		}
		err = json.Unmarshal([]byte(configData), &config)
		if err != nil {
			return nil, fmt.Errorf("Malformed config JSON!")
		}
	}

	options.Input = *input
	if options.Input == "" {
		return nil, fmt.Errorf("No input file given!")
	}

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

	options.Output = *output
	if options.Output == "" {
		options.Output = filepath.Join(filepath.Dir(options.Input), DEFAULT_FILENAME)
	}
	options.OutputDir = filepath.Dir(options.Output)

	options.CSP = *csp
	if options.CSP {
		dir, htmlFile := filepath.Split(options.Output)
		jsFile := htmlFile[:len(htmlFile)-len(".html")] + ".js"
		options.CSPFile = filepath.Join(dir, jsFile)
	}

	options.Inline = *inline
	options.Strip = *strip

	options.Verbose = *verbose
	options.Help = *help

	return options, nil
}
