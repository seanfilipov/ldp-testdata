package testdata

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/icrowley/fake"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type itemStatus struct {
	Name string `json:"name"`
}

type storageItem struct {
	ID                  string     `json:"id"`
	HoldingsRecordID    string     `json:"holdingsRecordId"`
	Barcode             string     `json:"barcode"`
	Status              itemStatus `json:"status"`
	Enumeration         string     `json:"enumeration"`
	CopyNumbers         []string   `json:"copyNumbers"`
	ItemLevelCallNumber string     `json:"itemLevelCallNumber"`
	PermanentLocationID string     `json:"permanentLocationId"`
	TemporaryLocationID string     `json:"temporaryLocationId"`
	MaterialTypeID      string     `json:"materialTypeID"`
}

// random returns a random int given a min/max range (include max)
func random(min, max int) int {
	return rand.Intn(1+max-min) + min
}
func randomEnumeration() string {
	randVolNum := random(1, 30)
	return fmt.Sprintf("v. %d", randVolNum)
}
func randomCopyNumbers() []string {
	return []string{strconv.Itoa(random(1, 5))}
}

func GenerateStorageItems(allParams AllParams, numItems int) {

	locChnl := streamRandomItem(allParams.Output, "locations.json", "locations")
	matChnl := streamRandomItem(allParams.Output, "material-types.json", "mtypes")

	rand.Seed(time.Now().UnixNano())
	makeStorageItem := func() storageItem {
		randomLocation, _ := <-locChnl
		randomMaterial, _ := <-matChnl
		var locationObj location
		var materialObj materialType
		mapstructure.Decode(randomLocation, &locationObj)
		mapstructure.Decode(randomMaterial, &materialObj)
		return storageItem{
			ID:                  uuid.Must(uuid.NewV4()).String(),
			HoldingsRecordID:    uuid.Must(uuid.NewV4()).String(),
			Barcode:             fake.DigitsN(16),
			Status:              itemStatus{Name: "Available"},
			Enumeration:         randomEnumeration(),
			CopyNumbers:         randomCopyNumbers(),
			ItemLevelCallNumber: randomCallNumber(),
			PermanentLocationID: locationObj.ID,
			MaterialTypeID:      materialObj.ID,
		}
	}
	var storageItems []interface{}
	for i := 0; i < numItems; i++ {
		oneStorageItem := makeStorageItem()
		storageItems = append(storageItems, oneStorageItem)
	}
	filename := "storageItems.json"
	objKey := "items"
	writeOutput(allParams.Output, filename, objKey, storageItems)

	updateManifest(FileDef{
		Module:    "mod-inventory-storage",
		Path:      "/item-storage/items",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  1,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/item-storage.html",
		N:         numItems,
	}, allParams.Output)
}
