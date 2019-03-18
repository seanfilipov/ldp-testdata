package testdata

// This purpose of this file is to retrieve a random LINE from a file (assumes line-by-line JSON)
// There are two functions to do this:
// 	1) read the whole file into memory, streamRandomLine()
// 	2) read the file line-by-line, streamRandomLineScan()
//
// Example usage:
//
// go streamRandomLine(filename, chnl)
// newVal, ok := <-chnl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// Returns the number of lines in the given file
func lineCounter(filename string) (int, error) {
	reader, err := os.Open(filename)
	defer reader.Close()
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 32*1024)
	count := 1
	lineSep := []byte{'\n'}
	for {
		c, err := reader.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// Scan file to line N and return it
func readLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	scanner := bufio.NewScanner(r)
	lineNum++
	for scanner.Scan() {
		lastLine++
		if lastLine == lineNum {
			return scanner.Text(), lastLine, scanner.Err()
		}
	}
	return line, lastLine, io.EOF
}

// Reads a whole file into memory and returns a slice of its lines
func readLines(filename string) ([]string, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	numLines := 0
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		numLines++
	}
	return lines, numLines, scanner.Err()
}

// Scan the file each time to stream back a random line
func streamRandomLineScan(filename string, chnl chan string) {
	maxValue, _ := lineCounter(filename) // Count the number of lines in the file
	for {
		jsonFile, _ := os.Open(filename)
		rand.Seed(time.Now().UnixNano())
		randomLineNumber := rand.Intn(maxValue)
		line, _, _ := readLine(jsonFile, randomLineNumber) // Read a random line from the file
		jsonFile.Close()
		chnl <- line
	}
	// close(chnl)
}

// Read the file into memory and stream back a random line
func streamRandomLine(filename string, chnl chan string) {
	maxValue, _ := lineCounter(filename) // Count the number of lines in the file
	lines, maxValue, _ := readLines(filename)
	for {
		rand.Seed(time.Now().UnixNano())
		randomLineNumber := rand.Intn(maxValue)
		line := lines[randomLineNumber]
		chnl <- line
		// TODO: break out of loop if chnl says so
	}
	// close(chnl)
}

// Parse the file as a JSON array, e.g. [{},{},...]
func streamRandomSliceItem(filename string, chnl chan interface{}) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []interface{}
	json.Unmarshal(byteValue, &result)
	for {
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(len(result))
		randomItem := result[randomNum]
		chnl <- randomItem
	}
}

// Parse the file as FOLIO JSON, e.g. {keyname:[{},{},...]}
func streamRandomFolioSliceItem(jsonKeyname, filename string, chnl chan interface{}) {
	jsonFile, errOpeningFile := os.Open(filename)
	if errOpeningFile != nil {
		logger.Error(errOpeningFile)
		chnl <- errOpeningFile
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)
	jsonArray, _ := result[jsonKeyname].([]interface{})
	for {
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(len(jsonArray))
		randomItem := jsonArray[randomNum]
		chnl <- randomItem
	}
}

// Linearly parse the file as a JSON array
func streamSliceItem(filename string, chnl chan interface{}) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []interface{}
	json.Unmarshal(byteValue, &result)
	for i := 0; i < len(result); i++ {
		chnl <- result[i]
	}
	close(chnl)
}

// Linearly parse the file as FOLIO JSON
func streamFolioSliceItem(jsonKeyname, filename string, chnl chan interface{}) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)
	jsonArray, _ := result[jsonKeyname].([]interface{})
	for i := 0; i < len(jsonArray); i++ {
		chnl <- jsonArray[i]
	}
	close(chnl)
}
