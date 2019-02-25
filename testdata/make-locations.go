package testdata

import (
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type location struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func GenerateLocations(outputParams OutputParams, numLocations int) {
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
	writeOutput(outputParams, "locations.json", "locations", locations)
}
