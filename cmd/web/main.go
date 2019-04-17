package main

import (
	"flag"
	"fmt"

	"github.com/folio-org/ldp-testdata/logging"
	"github.com/folio-org/ldp-testdata/testdata"
	"github.com/folio-org/ldp-testdata/web"
)

var logger = logging.Logger

func printUsage() {
	fmt.Println("\ngo run ./cmd/web/main.go [FLAGS]")
	fmt.Printf("\nAll flags are optional\n\n")
	flag.PrintDefaults() // Print the flag help strings
}

func main() {
	logging.Init()

	flag.Usage = func() {
		printUsage()
	}

	openBrowser := flag.Bool("openBrowser", true, "Whether to open a web browser to the UI")
	fileDefsFlag := flag.String("fileDefs", "filedefs.json", "The filepath of the JSON file definitions")
	flag.Parse()

	fileDefs := testdata.ParseFileDefs(*fileDefsFlag, "", false)
	web.Run(*openBrowser, fileDefs)
}
