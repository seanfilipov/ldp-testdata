package testdata

import (
	"path/filepath"

	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type location struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func GenerateLocations(outputDir string, numLocations int) {
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
	filepath := filepath.Join(outputDir, "locations.json")
	writeSliceToFile(filepath, locations, true)
}
