package testdata

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PrintUsage displays the usage ¯\_(ツ)_/¯
func PrintUsage() {
	fmt.Println("\ngo run ./cmd/ldp-testdata/main.go [FLAGS]")
	fmt.Printf("\nAll flags are optional\n\n")
	flag.PrintDefaults() // Print the flag help strings
}

// MakeTimestampedDir makes a timestamped directory if the -dir flag is unset
func MakeTimestampedDir(dirFlag string) string {
	if dirFlag != "" {
		os.MkdirAll(dirFlag, os.ModePerm) // Make the directory if it does not already exist
		return dirFlag
	}
	extractDir := "./output"
	currentTime := time.Now()
	timeStr := currentTime.Format("20060102_150405")
	outputDir := filepath.Join(extractDir, timeStr)
	os.MkdirAll(outputDir, os.ModePerm)
	return outputDir
}

// MapFileDefsToFunc checks that each 'path' is valid
func MapFileDefsToFunc(fileDefs []FileDef) (genFuncs []GenFunc) {
	validPaths := []string{
		"/groups",
		"/users",
		"/locations",
		"/material-types",
		"/instance-types",
		"/instance-storage/instances",
		"/holdings-storage/holdings",
		"/item-storage/items",
		"/inventory/items",
		"/loan-storage/loans",
		"/circulation/loans"}

	for _, def := range fileDefs {
		switch def.Path {
		case "/groups":
			genFuncs = append(genFuncs, GenerateGroups)
		case "/users":
			genFuncs = append(genFuncs, GenerateUsers)
		case "/locations":
			genFuncs = append(genFuncs, GenerateLocations)
		case "/material-types":
			genFuncs = append(genFuncs, GenerateMaterialTypes)
		case "/instance-types":
			genFuncs = append(genFuncs, GenerateInstanceTypes)
		case "/instance-storage/instances":
			genFuncs = append(genFuncs, GenerateInstances)
		case "/holdings-storage/holdings":
			genFuncs = append(genFuncs, GenerateHoldings)
		case "/item-storage/items":
			genFuncs = append(genFuncs, GenerateStorageItems)
		case "/inventory/items":
			genFuncs = append(genFuncs, GenerateInventoryItems)
		case "/loan-storage/loans":
			genFuncs = append(genFuncs, GenerateLoans)
		case "/circulation/loans":
			genFuncs = append(genFuncs, GenerateCirculationLoans)
		default:
			logger.Errorf("Error: '%s' is not a valid path value. \n  Valid paths: \n    %v\n\n", def.Path, strings.Join(validPaths, "\n    "))
			os.Exit(1)
		}
	}
	return
}

// ParseFileDefs reads filedefs on the command-line and/or the fileDefs.json file and returns a slice of fileDefs
func ParseFileDefs(filepath, fileDefsOverrideFlag string, onlyUseOverride bool) (fileDefs []FileDef) {
	var commandlineFileDefs []FileDef
	// 1) Parse the command-line filedefs, if any
	if fileDefsOverrideFlag != "" {
		marshalErr := json.Unmarshal([]byte(fileDefsOverrideFlag), &commandlineFileDefs)
		if marshalErr != nil {
			panic(marshalErr)
		}
	}
	// 2) Check if we're only using the command-line filedefs
	if onlyUseOverride {
		return commandlineFileDefs
	}

	// 3) Parse filedefs.json
	if filepath != "" {
		jsonFile, errOpenFile := os.Open(filepath)
		if errOpenFile != nil {
			fmt.Println("Error: Cannot find filedefs.json\nPlease run the command from the project root")
			os.Exit(1)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &fileDefs)
	}

	// 4) Merge the command-line values over the parsed values
	if fileDefsOverrideFlag != "" {
		// Merge in overrides
		for _, overrideDef := range commandlineFileDefs {
			for i := range fileDefs {
				if fileDefs[i].Path == overrideDef.Path {
					logger.Debugf("Using n=%d for %s\n", overrideDef.N, overrideDef.Path)
					fileDefs[i].N = overrideDef.N
					break
				}
			}
		}
	}
	return
}

func countFilesWithPrefix(filepath, prefix string) (numMatching int) {
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		logger.Fatal(err)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) {
			numMatching++
			// fmt.Println(f.Name())
		}
	}
	return numMatching
}
