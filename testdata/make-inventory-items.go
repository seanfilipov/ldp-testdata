package testdata

// UNUSED
// This file is unused because mod-inventory is a business logic module
// instead of the storage module (mod-loan-storage)
// https://s3.amazonaws.com/foliodocs/api/mod-inventory/inventory.html

// type materialType struct {
// 	Name string `json:"name"`
// }
// type itemStatus struct {
// 	Name string `json:"name"`
// }
// type inventoryItem struct {
// 	Title             string       `json:"title"`
// 	ID                string       `json:"id"`
// 	Barcode           string       `json:"barcode"`
// 	HoldingsRecordID  string       `json:"holdingsRecordId"`
// 	EffectiveLocation location     `json:"effectiveLocation"`
// 	Status            itemStatus   `json:"status"`
// 	MaterialType      materialType `json:"materialType"`
// }

// itemsChnl := streamOutputLinearly(outputParams, "items.json", "items")
// for oneItem := range itemsChnl {
// 	var itemObj item
// 	mapstructure.Decode(oneItem, &itemObj)
// 	oneStorageItem := storageItem{
// 		ID:               itemObj.ID,
// 		HoldingsRecordID: itemObj.HoldingsRecordID,
// 		Barcode:          itemObj.Barcode,
// 		Status:           itemObj.Status,

// func GenerateInventoryItems(outputParams OutputParams) {
// 	bookChnl := make(chan string, 1)
// 	_, nameOfThisFile, _, _ := runtime.Caller(0)
// 	pkgDir := path.Dir(nameOfThisFile)
// 	go streamRandomLine(pkgDir+"/book_titles.txt", bookChnl)

// 	locChnl := streamRandomItem(outputParams, "locations.json", "locations")
// 	makeItem := func() item {
// 		// TODO: Should iterate over titles, not get a random one
// 		randomBookTitle, _ := <-bookChnl
// 		randomLocation, _ := <-locChnl
// 		var locationObj location
// 		mapstructure.Decode(randomLocation, &locationObj)
// 		return item{
// 			Title:             randomBookTitle,
// 			ID:                uuid.Must(uuid.NewV4()).String(),
// 			Barcode:           fake.DigitsN(16),
// 			HoldingsRecordID:  uuid.Must(uuid.NewV4()).String(),
// 			EffectiveLocation: locationObj,
// 			Status:            itemStatus{Name: "Available"},
// 			MaterialType:      materialType{Name: "book"},
// 		}
// 	}
// 	var items []interface{}
// 	itemsChnl := streamOutputLinearly(outputParams, "items.json", "items")
// 	for oneItem := range itemsChnl {
// 		var itemObj storageItem
// 		mapstructure.Decode(oneItem, &itemObj)
// 		u := makeItem()
// 		items = append(items, u)
// 	}
// 	writeOutput(outputParams, "items.json", "items", items)
// }
