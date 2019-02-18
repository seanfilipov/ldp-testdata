package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

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

// StorageItems share the same information with inventory items
// (ID, holdingsRecordId, barcode)
func generateStorageItems(outputDir string) {
	rand.Seed(time.Now().UnixNano())
	var storageItems []interface{}
	itemsPath := filepath.Join(outputDir, "items.json")
	itemsChnl := make(chan interface{}, 1)
	go streamSliceItem(itemsPath, itemsChnl)
	for oneItem := range itemsChnl {
		var itemObj item
		mapstructure.Decode(oneItem, &itemObj)
		oneStorageItem := storageItem{
			ID:               itemObj.ID,
			HoldingsRecordID: itemObj.HoldingsRecordID,
			Barcode:          itemObj.Barcode,
			Status:           itemObj.Status,
			Enumeration:      randomEnumeration(),
			CopyNumbers:      randomCopyNumbers(),
		}
		storageItems = append(storageItems, oneStorageItem)
	}
	filepath := filepath.Join(outputDir, "storageItems.json")
	writeSliceToFile(filepath, storageItems, true)
}
