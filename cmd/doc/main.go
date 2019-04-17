package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/folio-org/ldp-testdata/logging"
	"github.com/folio-org/ldp-testdata/testdata"
)

var logger = logging.Logger

func printUsage() {
	fmt.Println("\ngo run ./cmd/doc/main.go [FLAGS]")
	fmt.Printf("\nAll flags are optional\n\n")
	flag.PrintDefaults() // Print the flag help strings
}

// This script auto-updates the README's Supported Routes section
// based on filedefs.json
func main() {
	logging.Init()

	flag.Usage = func() {
		printUsage()
	}

	fileDefsFlag := flag.String("fileDefs", "filedefs.json", "The filepath of the JSON file definitions")
	flag.Parse()
	fileDefs := testdata.ParseFileDefs(*fileDefsFlag, "", false)

	file, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	atSupportedRoutes := false
	for scanner.Scan() {
		text := scanner.Text()
		if text == "Supported Routes" {
			atSupportedRoutes = true
		} else if atSupportedRoutes && strings.HasPrefix(text, "- ") {
			for strings.HasPrefix(scanner.Text(), "- ") {
				scanner.Scan()
			}
			for _, fileDef := range fileDefs {
				newText := fmt.Sprintf("- [%s](%s)", fileDef.Path, fileDef.Doc)
				lines = append(lines, newText)
			}
			atSupportedRoutes = false
			continue
		}
		lines = append(lines, text)
	}
	writeSliceToFileLineByLine("README.md", lines)
}

// Writes a file line by line
func writeSliceToFileLineByLine(filepath string, slice []string) {
	f, _ := os.Create(filepath)
	defer f.Close()
	w := bufio.NewWriter(f)
	for i := 0; i < len(slice); i++ {
		w.WriteString(slice[i] + "\n")
	}
	w.Flush()
}
