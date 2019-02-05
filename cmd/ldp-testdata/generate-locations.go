package main

import (
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type location struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func generateLocations(filepath string) {
	makeLocation := func() location {
		return location{
			Name: fake.LastName() + " Library",
			ID:   uuid.Must(uuid.NewV4()).String(),
		}
	}
	var locations []interface{}
	for i := 0; i < 20; i++ {
		l := makeLocation()
		locations = append(locations, l)
	}
	writeSliceToFile(filepath, locations, true)
}
