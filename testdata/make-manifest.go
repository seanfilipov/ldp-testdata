package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileDef struct {
	Module    string `json:"module"`            // the module
	Path      string `json:"path"`              // API route simulated
	ObjectKey string `json:"objectKey"`         // the field that contains the array in the output JSON
	NumFiles  int    `json:"numFiles"`          // the number of files a part of this output
	Doc       string `json:"doc"`               // URL to the API documentation
	N         int    `json:"n,omitempty"`       // Number of objects
	NumDays   int    `json:"numDays,omitempty"` // (Optional) Time interval to spread the objects across
}

func toInterface(originals []FileDef) []interface{} {
	newThings := make([]interface{}, len(originals))
	for i, s := range originals {
		newThings[i] = s
	}
	return newThings
}

func updateManifest(def FileDef, params OutputParams) {
	filepath := filepath.Join(params.OutputDir, "manifest.json")

	jsonFile, errOpeningFile := os.Open(filepath)
	if errOpeningFile != nil {
		// write file
		var fileDefs []interface{}
		fileDefs = append(fileDefs, def)
		writeSliceToFile(filepath, fileDefs, true)
	} else {
		// read JSON, then update it

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}
		var defs []FileDef
		json.Unmarshal(byteValue, &defs)
		foundTarget := false
		for i := 0; i < len(defs); i++ {
			if defs[i].Path == def.Path {
				logger.Debugf("Overwriting entry for %s", def.Path)
				defs[i] = def
				foundTarget = true
				break
			}
		}
		if !foundTarget {
			defs = append(defs, def)
		}
		newDefs := toInterface(defs)
		writeSliceToFile(filepath, newDefs, true)
	}
}
