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

func GenerateLocations(allParams AllParams, numLocations int) {
	makeLocation := func() location {
		return location{
			Name: fake.LastName() + " Library",
			ID:   uuid.Must(uuid.NewV4()).String(),
		}
	}
	var locations []interface{}
	for i := 0; i < numLocations; i++ {
		l := makeLocation()
		locations = append(locations, l)
	}

	filename := "locations.json"
	objKey := "locations"
	writeOutput(allParams.Output, filename, objKey, locations)

	updateManifest(FileDef{
		Module:    "mod-inventory-storage",
		Path:      "/locations",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  1,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/location.html",
		N:         numLocations,
	}, allParams.Output)
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
