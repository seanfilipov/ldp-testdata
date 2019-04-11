package main

import (
	"flag"

	"github.com/folio-org/ldp-testdata/logging"
	"github.com/folio-org/ldp-testdata/testdata"
)

var logger = logging.Logger

func main() {
	logging.Init()
	flag.Usage = func() {
		testdata.PrintUsage()
	}
	// openBrowser := flag.Bool("openBrowser", true, "Whether to open a web browser to the UI")
	dirFlag := flag.String("dir", "", "The directory to store output")
	fileDefsFlag := flag.String("fileDefs", "filedefs.json", "The filepath of the JSON file definitions")
	dataFmtFlag := flag.String("dataFormat", "folioJSON", `The outputted data format [folioJSON|jsonArray]`)
	fileDefsOverrideFlag := flag.String("json", "", `JSON array to override the number of objects set filedefs.json
Example: '[{"path": "/loan-storage/loans", "n":50000}]'`)
	onlyUseOverrideFlag := flag.Bool("only-json", false, "Use with the -json flag to ignore filedefs.json")
	flag.Parse()

	fileDefs := testdata.ParseFileDefs(*fileDefsFlag,
		*fileDefsOverrideFlag,
		*onlyUseOverrideFlag)
	funcs := testdata.MapFileDefsToFunc(fileDefs)

	// web.Run(*openBrowser, fileDefs)

	// If we need to do any more validation of params, change this to a NewParams() function
	// which does the validation
	p := testdata.AllParams{
		FileDefs: fileDefs,
		Output: testdata.OutputParams{
			OutputDir:  testdata.MakeTimestampedDir(*dirFlag),
			DataFormat: testdata.ParseDataFmtFlag(*dataFmtFlag),
			Indent:     true,
		},
	}
	testdata.MakeAll(funcs, p)
	logger.Infof("Generated data in %s\n", p.Output.OutputDir)
}
