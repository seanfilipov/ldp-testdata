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

func GenerateStorageItems(filedef FileDef, outputParams OutputParams) {

	holdingChnl := streamOutputLinearly(outputParams, "holdings-storage-holdings-1.json", "holdingsRecords")
	locChnl := streamRandomItem(outputParams, "locations-1.json", "locations")
	matChnl := streamRandomItem(outputParams, "material-types-1.json", "mtypes")

	rand.Seed(time.Now().UnixNano())
	makeStorageItem := func(oneHolding interface{}) storageItem {
		randomLocation, _ := <-locChnl
		randomMaterial, _ := <-matChnl
		var locationObj location
		var materialObj materialType
		var holdingObj holding
		mapstructure.Decode(randomLocation, &locationObj)
		mapstructure.Decode(randomMaterial, &materialObj)
		mapstructure.Decode(oneHolding, &holdingObj)
		return storageItem{
			ID:                  uuid.Must(uuid.NewV4()).String(),
			HoldingsRecordID:    holdingObj.ID,
			Barcode:             fake.DigitsN(16),
			Status:              itemStatus{Name: "Available"},
			Enumeration:         randomEnumeration(),
			CopyNumbers:         randomCopyNumbers(),
			ItemLevelCallNumber: holdingObj.CallNumber,
			PermanentLocationID: locationObj.ID,
			MaterialTypeID:      materialObj.ID,
		}
	}
	var storageItems []interface{}
	for oneHolding := range holdingChnl {
		oneStorageItem := makeStorageItem(oneHolding)
		storageItems = append(storageItems, oneStorageItem)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, storageItems)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
