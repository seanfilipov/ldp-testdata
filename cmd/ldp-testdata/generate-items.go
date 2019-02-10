package main

import (
	"path/filepath"

	"github.com/icrowley/fake"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type materialType struct {
	Name string `json:"name"`
}
type itemStatus struct {
	Name string `json:"name"`
}
type item struct {
	Title             string       `json:"title"`
	ID                string       `json:"id"`
	Barcode           string       `json:"barcode"`
	HoldingsRecordID  string       `json:"holdingsRecordId"`
	EffectiveLocation location     `json:"effectiveLocation"`
	Status            itemStatus   `json:"status"`
	MaterialType      materialType `json:"materialType"`
}

func generateItems(outputDir string) {
	bookChnl := make(chan string, 1)
	locChnl := make(chan interface{}, 1)
	locationsPath := filepath.Join(outputDir, "locations.json")
	go streamRandomSliceItem(locationsPath, locChnl)
	go streamRandomLine("./book_titles.txt", bookChnl)
	makeItem := func() item {
		// TODO: Should iterate over titles, not get a random one
		randomBookTitle, _ := <-bookChnl
		randomLocation, _ := <-locChnl
		var locationObj location
		mapstructure.Decode(randomLocation, &locationObj)
		return item{
			Title:             randomBookTitle,
			ID:                uuid.Must(uuid.NewV4()).String(),
			Barcode:           fake.DigitsN(16),
			HoldingsRecordID:  uuid.Must(uuid.NewV4()).String(),
			EffectiveLocation: locationObj,
			Status:            itemStatus{Name: "Available"},
			MaterialType:      materialType{Name: "book"},
		}
	}
	var items []interface{}
	for i := 0; i < 9980; i++ {
		u := makeItem()
		items = append(items, u)
	}
	filepath := filepath.Join(outputDir, "items.json")
	writeSliceToFile(filepath, items, true)
}
