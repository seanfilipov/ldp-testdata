package testdata

import "path/filepath"

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
	Output       OutputParams
	NumGroups    int
	NumUsers     int
	NumLocations int
	NumItems     int
	NumLoans     int
}

func MakeAll(p AllParams) {
	GenerateGroups(p.Output, p.NumGroups)
	GenerateUsers(p.Output, p.NumUsers)
	GenerateLocations(p.Output, p.NumLocations)
	GenerateStorageItems(p.Output, p.NumItems)
	GenerateLoans(p.Output, p.NumLoans)
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
