package testdata

import (
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type location struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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
	}, allParams.Output)
}
