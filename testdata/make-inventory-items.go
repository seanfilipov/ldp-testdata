package testdata

import (
	"path"
	"runtime"

	"github.com/mitchellh/mapstructure"
)

// This file is deprecated because mod-inventory is a business logic module
// instead of the storage module (mod-loan-storage)
// https://s3.amazonaws.com/foliodocs/api/mod-inventory/inventory.html

type materialType struct {
	Name string `json:"name"`
}
type inventoryItem struct {
	Title             string       `json:"title"`
	ID                string       `json:"id"`
	Barcode           string       `json:"barcode"`
	HoldingsRecordID  string       `json:"holdingsRecordId"`
	EffectiveLocation location     `json:"effectiveLocation"`
	Status            itemStatus   `json:"status"`
	MaterialType      materialType `json:"materialType"`
}

func GenerateInventoryItems(allParams AllParams, numItems int) {
	outputParams := allParams.Output
	bookChnl := make(chan string, 1)
	_, nameOfThisFile, _, _ := runtime.Caller(0)
	pkgDir := path.Dir(nameOfThisFile)
	go streamRandomLine(pkgDir+"/book_titles.txt", bookChnl)

	locChnl := streamRandomItem(outputParams, "locations.json", "locations")
	makeItem := func(storageItemObj storageItem) inventoryItem {
		// TODO: Should iterate over titles, not get a random one
		randomBookTitle, _ := <-bookChnl
		randomLocation, _ := <-locChnl
		var locationObj location
		mapstructure.Decode(randomLocation, &locationObj)
		return inventoryItem{
			Title:             randomBookTitle,
			ID:                storageItemObj.ID,
			Barcode:           storageItemObj.Barcode,
			HoldingsRecordID:  storageItemObj.HoldingsRecordID,
			EffectiveLocation: locationObj,
			Status:            storageItemObj.Status,
			MaterialType:      materialType{Name: "book"},
		}
	}
	var items []interface{}
	itemsChnl := streamOutputLinearly(outputParams, "storageItems.json", "items")
	for oneItem := range itemsChnl {
		var storageItemObj storageItem
		mapstructure.Decode(oneItem, &storageItemObj)
		u := makeItem(storageItemObj)
		items = append(items, u)
	}
	filename := "inventoryItems.json"
	objKey := "items"
	writeOutput(outputParams, filename, objKey, items)

	updateManifest(FileDef{
		Module:    "mod-inventory",
		Path:      "/inventory/items",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  1,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-inventory/inventory.html",
		N:         len(items),
	}, allParams.Output)
}
