package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type location struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type locationsFile struct {
	Locations []location `json:"locations"`
}

func GenerateLocations(filedef FileDef, outputParams OutputParams) {
	makeLocation := func() location {
		return location{
			Name: fake.LastName() + " Library",
			ID:   uuid.Must(uuid.NewV4()).String(),
		}
	}
	var locations []interface{}
	for i := 0; i < filedef.N; i++ {
		l := makeLocation()
		locations = append(locations, l)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, locations)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}

//
// Helpers for other files:
//

func readLocations(params OutputParams, filename string) []location {
	filepath := filepath.Join(params.OutputDir, filename)
	jsonFile, errOpeningFile := os.Open(filepath)
	if errOpeningFile != nil {
		panic(errOpeningFile)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	var locationsFileObj locationsFile
	json.Unmarshal(byteValue, &locationsFileObj)
	return locationsFileObj.Locations
}

func lookupLocation(ID string, locations *[]location) location {
	var matchingLoc location
	for _, loc := range *locations {
		if loc.ID == ID {
			matchingLoc = loc
			break
		}
	}
	return matchingLoc
}
