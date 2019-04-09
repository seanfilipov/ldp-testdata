package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/folio-org/ldp-testdata/logging"
	"github.com/folio-org/ldp-testdata/testdata"
	"github.com/folio-org/ldp-testdata/web"
)

var logger = logging.Logger

func makeTimestampedDir(dirFlag string) string {
	if dirFlag != "" {
		os.MkdirAll(dirFlag, os.ModePerm) // Make the directory if it does not already exist
		return dirFlag
	}
	extractDir := "./extract-output"
	currentTime := time.Now()
	timeStr := currentTime.Format("20060102_150405")
	outputDir := filepath.Join(extractDir, timeStr)
	os.MkdirAll(outputDir, os.ModePerm)
	return outputDir
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("./ldp-testdata FLAGS [all|groups|users|locations|loans|storageitems]")
	fmt.Println("  where FLAGS include:")
	flag.PrintDefaults() // Print the flag help strings
}

func main() {
	logging.Init()
	flag.Usage = func() {
		printUsage()
	}
	openBrowser := flag.Bool("openBrowser", true, "Whether to open a web browser to the UI")
	dirFlag := flag.String("dir", "",
		`The directory to use for extract output. If the selected test data depends on
other test data (e.g. 'users' depends on 'groups'), that dependency should exist
in this directory.`)
	fileDefsFlag := flag.String("fileDefs", "", "The filepath of the JSON file definitions")
	dataFmtFlag := flag.String("dataFormat", "folioJSON", `The outputted data format [folioJSON|jsonArray]`)
	numGroupsFlag := flag.Int("nGroups", 12, `The number of groups to create`)
	numUsersFlag := flag.Int("nUsers", 30000, `The number of users to create`)
	numLocationsFlag := flag.Int("nLocations", 20, `The number of locations to create`)
	numItemsFlag := flag.Int("nItems", 10000, `The number of items to create`)
	numLoansFlag := flag.Int("nLoans", 10000, `The number of loans to create`)
	flag.Parse()

	// VALIDATE ARGUMENT 'MODE' IS VALID
	modes := map[string]bool{
		"all":          true,
		"groups":       true,
		"users":        true,
		"locations":    true,
		"loans":        true,
		"storageitems": true,
	}
	fileDefs := testdata.ParseFileDefs(*fileDefsFlag)
	web.Run(*openBrowser, fileDefs)

	if len(flag.Args()) < 1 {
		printUsage()
		os.Exit(1)
	}
	mode := flag.Arg(0)
	if _, ok := modes[mode]; !ok {
		fmt.Printf("Error: '%s' is not a valid argument\n", mode)
		printUsage()
		os.Exit(1)
	}
	// If we need to do any more validation of params, change this to a NewParams() function
	// which does the validation
	p := testdata.AllParams{
		FileDefs: fileDefs,
		Output: testdata.OutputParams{
			OutputDir:  makeTimestampedDir(*dirFlag),
			DataFormat: testdata.ParseDataFmtFlag(*dataFmtFlag),
			Indent:     true,
		},
		NumGroups:    *numGroupsFlag,
		NumUsers:     *numUsersFlag,
		NumLocations: *numLocationsFlag,
		NumItems:     *numItemsFlag,
		NumLoans:     *numLoansFlag,
	}

	switch mode {
	case "all":
		testdata.MakeAll(p)
	case "groups":
		testdata.GenerateGroups(p, p.NumGroups)
	case "users":
		testdata.GenerateUsers(p, p.NumUsers)
	case "locations":
		testdata.GenerateLocations(p, p.NumLocations)
	case "storageitems":
		testdata.GenerateStorageItems(p, p.NumItems)
	case "loans":
		testdata.GenerateLoans(p, p.NumLoans)
	}
	logger.Infof("Generated data in %s\n", p.Output.OutputDir)
}
