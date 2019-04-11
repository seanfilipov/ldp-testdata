package testdata

import (
	"path/filepath"
)

// DataFmt is an enum type
type DataFmt int

const (
	// FolioJSON is the default data format output, the FOLIO JSON format used today
	FolioJSON DataFmt = iota
	// JSONArray is the same as FolioJSON but without the key that wraps the array
	JSONArray
)

type OutputParams struct {
	OutputDir  string
	DataFormat DataFmt
	Indent     bool
}

// AllParams includes the input parameters FileDefs, and the OutputParams
type AllParams struct {
	FileDefs []FileDef
	Output   OutputParams
}

// GenFunc is a function that generates a data output
type GenFunc func(AllParams, int)

func MakeAll(funcs []GenFunc, p AllParams) {
	for i, fileDef := range p.FileDefs {
		funcs[i](p, fileDef.N)
	}
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
