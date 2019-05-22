package testdata

import (
	"path"
	"runtime"

	"github.com/mitchellh/mapstructure"
)

// This file is deprecated because mod-inventory is a business logic module
// instead of the storage module (mod-loan-storage)
// https://s3.amazonaws.com/foliodocs/api/mod-inventory/inventory.html

type itemMaterialType struct {
	Name string `json:"name"`
}
type inventoryItem struct {
	Title             string           `json:"title"`
	ID                string           `json:"id"`
	Barcode           string           `json:"barcode"`
	HoldingsRecordID  string           `json:"holdingsRecordId"`
	EffectiveLocation location         `json:"effectiveLocation"`
	Status            itemStatus       `json:"status"`
	MaterialType      itemMaterialType `json:"materialType"`
}

func GenerateInventoryItems(filedef FileDef, outputParams OutputParams) {
	bookChnl := make(chan string, 1)
	_, nameOfThisFile, _, _ := runtime.Caller(0)
	pkgDir := path.Dir(nameOfThisFile)
	go streamRandomLine(pkgDir+"/book_titles.txt", bookChnl)

	locations := readLocations(outputParams, "locations-1.json")

	makeItem := func(storageItemObj storageItem) inventoryItem {
		// TODO: Should iterate over titles, not get a random one
		randomBookTitle, _ := <-bookChnl
		effectiveLocation := lookupLocation(storageItemObj.PermanentLocationID, &locations)
		return inventoryItem{
			Title:             randomBookTitle,
			ID:                storageItemObj.ID,
			Barcode:           storageItemObj.Barcode,
			HoldingsRecordID:  storageItemObj.HoldingsRecordID,
			EffectiveLocation: effectiveLocation,
			Status:            storageItemObj.Status,
			MaterialType:      itemMaterialType{Name: "book"},
		}
	}
	var items []interface{}
	itemsChnl := streamOutputLinearly(outputParams, "item-storage-items-1.json", "items")
	for oneItem := range itemsChnl {
		var storageItemObj storageItem
		mapstructure.Decode(oneItem, &storageItemObj)
		u := makeItem(storageItemObj)
		items = append(items, u)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, items)
	filedef.NumFiles = 1
	filedef.N = len(items)
	updateManifest(filedef, outputParams)
}
