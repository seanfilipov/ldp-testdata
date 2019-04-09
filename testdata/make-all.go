package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DataFmt int

const (
	FolioJSON DataFmt = iota
	JSONArray
)

type OutputParams struct {
	OutputDir  string
	DataFormat DataFmt
	Indent     bool
}
type AllParams struct {
	FileDefs     []FileDef
	Output       OutputParams
	NumGroups    int
	NumUsers     int
	NumLocations int
	NumItems     int
	NumLoans     int
}

// ParseFileDefs reads the fileDefs.json file and returns a slice of fileDefs
func ParseFileDefs(filepath string) (fileDefs []FileDef) {
	if filepath != "" {
		jsonFile, errOpenFile := os.Open(filepath)
		if errOpenFile != nil {
			panic(errOpenFile)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &fileDefs)
	}
	return
}

func MakeAll(p AllParams) {
	GenerateGroups(p, p.NumGroups)
	GenerateUsers(p, p.NumUsers)
	GenerateLocations(p, p.NumLocations)
	GenerateStorageItems(p, p.NumItems)
	GenerateLoans(p, p.NumLoans)
}

func writeOutput(params OutputParams, filename, jsonKeyname string, slice []interface{}) {
	filepath := filepath.Join(params.OutputDir, filename)
	if params.DataFormat == JSONArray {
		writeSliceToFile(filepath, slice, params.Indent)
	} else {
		writeFolioSliceToFile(jsonKeyname, filepath, slice, params.Indent)
	}
}

func streamRandomItem(params OutputParams, filename, jsonKeyname string) chan interface{} {
	chnl := make(chan interface{}, 1)
	filepath := filepath.Join(params.OutputDir, filename)
	if params.DataFormat == JSONArray {
		go streamRandomSliceItem(filepath, chnl)
	} else {
		go streamRandomFolioSliceItem(jsonKeyname, filepath, chnl)
	}
	return chnl
}
func streamOutputLinearly(params OutputParams, filename, jsonKeyname string) chan interface{} {
	chnl := make(chan interface{}, 1)
	filepath := filepath.Join(params.OutputDir, filename)
	if params.DataFormat == JSONArray {
		go streamSliceItem(filepath, chnl)
	} else {
		go streamFolioSliceItem(jsonKeyname, filepath, chnl)
	}
	return chnl
}

func ParseDataFmtFlag(flag string) DataFmt {
	if flag == "folioJSON" {
		return FolioJSON
	}
	return JSONArray // otherwise
}
