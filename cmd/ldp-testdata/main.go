package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/folio-org/ldp-testdata/testdata"
)

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
	fmt.Println("./ldp-testdata FLAGS [all|groups|users|locations|items|loans|circloans|storageitems]")
	fmt.Println("  where FLAGS include:")
	flag.PrintDefaults() // Print the flag help strings
}

func main() {
	flag.Usage = func() {
		printUsage()
	}
	dirFlag := flag.String("dir", "",
		`The directory to use for extract output. If the selected test data depends on
other test data (e.g. 'users' depends on 'groups'), that dependency should exist
in this directory.`)
	numGroupsFlag := flag.Int("nGroups", 12, `The number of groups to create`)
	numUsersFlag := flag.Int("nUsers", 30000, `The number of users to create`)
	numLocationsFlag := flag.Int("nLocations", 20, `The number of locations to create`)
	numLoansFlag := flag.Int("nLoans", 10000, `The number of loans to create`)
	flag.Parse()
	if len(flag.Args()) >= 1 {
		mode := flag.Arg(0)
		switch mode {
		case "all":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateGroups(timestampedDir, *numGroupsFlag)
			testdata.GenerateUsers(timestampedDir, *numUsersFlag)
			testdata.GenerateLocations(timestampedDir, *numLocationsFlag)
			testdata.GenerateItems(timestampedDir)
			testdata.GenerateLoans(timestampedDir, *numLoansFlag)
			testdata.GenerateCirculationLoans(timestampedDir)
			testdata.GenerateStorageItems(timestampedDir)
			fmt.Printf("Generated data in %s\n", timestampedDir)

		case "groups":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateGroups(timestampedDir, *numGroupsFlag)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		case "users":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateUsers(timestampedDir, *numUsersFlag)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		case "locations":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateLocations(timestampedDir, *numLocationsFlag)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		case "items":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateItems(timestampedDir)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		case "loans":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateLoans(timestampedDir, *numLoansFlag)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		case "circloans":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateCirculationLoans(timestampedDir)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		case "storageitems":
			timestampedDir := makeTimestampedDir(*dirFlag)
			testdata.GenerateStorageItems(timestampedDir)
			fmt.Printf("Generated data in %s\n", timestampedDir)
		default:
			fmt.Printf("Error: '%s' is not a valid argument\n", mode)
			printUsage()
		}
	} else {
		printUsage()
	}

	// threeGenerators()
}
