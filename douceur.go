package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/slt/douceur/inliner"
	"github.com/slt/douceur/parser"
)

const (
	// Version is package version
	Version = "0.3.0"
)

var (
	flagVersion bool
	cssPath     string
)

func init() {
	flag.BoolVar(&flagVersion, "version", false, "Display version")
	flag.StringVar(&cssPath, "es", "", "Include external stylesheet when inlining")
}

func main() {
	flag.Parse()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	action := flag.Arg(0)

	if action == "" {
		fmt.Println("No action supplied")
		usage()
		os.Exit(1)
	}

	switch action {
	case "parse":
		parseCSS(flag.Arg(1))
	case "inline":
		inlineCSS(flag.Arg(1), cssPath)
	default:
		fmt.Println("Unexpected action: ", action)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Usage: %s %s %s\n", os.Args[0], "(parse|inline)", "/path/to/file")
	fmt.Printf("Usage: %s %s < %s\n", os.Args[0], "(parse|inline)", "/path/to/file")
}

// parse and display CSS file
func parseCSS(filePath string) {
	input := read(filePath)

	stylesheet, err := parser.Parse(string(input))
	if err != nil {
		fmt.Println("Parsing error: ", err)
		os.Exit(1)
	}

	fmt.Println(stylesheet.String())
}

// inlines CSS into HTML and display result
func inlineCSS(filePath string, cssPath string) {
	htmlInput := string(read(filePath))
	cssInput := ""

	if cssPath != "" {
		cssInput = string(readFile(cssPath))
	}

	output, err := inliner.InlineWithExternalCSS(htmlInput, cssInput)
	if err != nil {
		fmt.Println("Inlining error: ", err)
		os.Exit(1)
	}

	fmt.Println(output)
}

func read(filePath string) []byte {
	fmt.Println(filePath)
	if filePath == "" {
		return readStandardInput()
	}
	return readFile(filePath)
}

func readFile(filePath string) []byte {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to open file: ", filePath, err)
		os.Exit(1)
	}

	return file
}

func readStandardInput() []byte {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Failed read stdin: ", err)
		os.Exit(1)
	}

	return data
}
