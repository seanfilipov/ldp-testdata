package main

// The purpose of this file is to write Go slices to file in the JSON format.
// There are two functions to do this:
// 	1) write a valid JSON array to file, writeSliceToFile()
// 	2) write 1 JSON object per line, writeSliceToFileLineByLine()
//
// The second method is available in case you want to read the JSON array line-by-line,
// instead of reading the whole file at one time.

import (
	"bufio"
	"encoding/json"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Writes valid JSON to a single line of a file
func writeSliceToFile(filepath string, jsonSlice []interface{}, indent bool) {
	// Create a file
	f, err := os.Create(filepath)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	// Write JSON to file
	if indent {
		byteJSON, _ := json.MarshalIndent(jsonSlice, "", "    ")
		w.Write(byteJSON)
	} else {
		byteJSON, _ := json.Marshal(jsonSlice)
		w.Write(byteJSON)
	}
	w.WriteString("\n")
	w.Flush()
}

// Writes one JSON object per line
func writeSliceToFileLineByLine(filepath string, jsonSlice []interface{}) {
	f, err := os.Create(filepath)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	// Write each JSON object as a line in the file
	for i := 0; i < len(jsonSlice); i++ {
		byteGroup, _ := json.Marshal(jsonSlice[i])
		w.Write(byteGroup)
		w.WriteString("\n")
	}
	w.Flush()
}
