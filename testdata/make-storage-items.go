package testdata

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type itemStatus struct {
	Name string `json:"name"`
}

type storageItem struct {
	ID               string     `json:"id"`
	HoldingsRecordID string     `json:"holdingsRecordId"`
	Barcode          string     `json:"barcode"`
	Status           itemStatus `json:"status"`
	Enumeration      string     `json:"enumeration"`
	CopyNumbers      []string   `json:"copyNumbers"`
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
func randomEnumeration() string {
	randVolNum := random(1, 30)
	randYear := random(1945, 2019)
	return fmt.Sprintf("v. %d %d", randVolNum, randYear)
}
func randomCopyNumbers() []string {
	return []string{strconv.Itoa(random(1, 5))}
}

func GenerateStorageItems(outputParams OutputParams, numItems int) {
	rand.Seed(time.Now().UnixNano())
	makeStorageItem := func() storageItem {
		return storageItem{
			ID:               uuid.Must(uuid.NewV4()).String(),
			HoldingsRecordID: uuid.Must(uuid.NewV4()).String(),
			Barcode:          fake.DigitsN(16),
			Status:           itemStatus{Name: "Available"},
			Enumeration:      randomEnumeration(),
			CopyNumbers:      randomCopyNumbers(),
		}
	}
	var storageItems []interface{}
	for i := 0; i < numItems; i++ {
		oneStorageItem := makeStorageItem()
		storageItems = append(storageItems, oneStorageItem)
	}
	writeOutput(outputParams, "storageItems.json", "items", storageItems)
}
