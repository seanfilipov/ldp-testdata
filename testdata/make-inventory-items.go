package testdata

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// https://github.com/folio-org/mod-inventory/blob/master/ramls/item.json
// https://s3.amazonaws.com/foliodocs/api/mod-inventory/inventory.html

// FYI this path is for a business logic module (mod-inventory)
// instead of the storage module (mod-inventory-storage)

type itemMaterialType struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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

type holdingsFile struct {
	Holdings []holding `json:"holdingsRecords"`
}

func readHoldings(params OutputParams, filename string) []holding {
	byteValue := readFile(params, filename)
	var holdingsFileObj holdingsFile
	json.Unmarshal(byteValue, &holdingsFileObj)
	return holdingsFileObj.Holdings
}
func getHoldingsMap(params OutputParams, filename string) map[string]string {
	holdingsMap := make(map[string]string)
	holdings := readHoldings(params, filename)
	for _, elmt := range holdings {
		holdingsMap[elmt.ID] = elmt.ShelvingTitle
	}
	return holdingsMap
}

func GenerateInventoryItems(filedef FileDef, outputParams OutputParams) {

	holdings := getHoldingsMap(outputParams, "holdings-storage-holdings-1.json")
	locations := readLocations(outputParams, "locations-1.json")
	matChnl := streamRandomItem(outputParams, "material-types-1.json", "mtypes")
	makeItem := func(storageItemObj storageItem) inventoryItem {
		// TODO: Should iterate over titles, not get a random one
		randomMaterial, _ := <-matChnl
		var materialObj materialType
		mapstructure.Decode(randomMaterial, &materialObj)

		effectiveLocation := lookupLocation(storageItemObj.PermanentLocationID, &locations)
		return inventoryItem{
			Title:             holdings[storageItemObj.HoldingsRecordID],
			ID:                storageItemObj.ID,
			Barcode:           storageItemObj.Barcode,
			HoldingsRecordID:  storageItemObj.HoldingsRecordID,
			EffectiveLocation: effectiveLocation,
			Status:            storageItemObj.Status,
			MaterialType:      itemMaterialType{ID: materialObj.ID, Name: materialObj.Name},
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
