package testdata

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type circulationLoanItem struct {
	Title            string     `json:"title"`
	Barcode          string     `json:"barcode"`
	HoldingsRecordID string     `json:"holdingsRecordId"`
	Status           itemStatus `json:"status"`
	Location         itemStatus `json:"location"`
}

type circulationLoan struct {
	ID       string              `json:"id"`
	UserID   string              `json:"userId"`
	ItemID   string              `json:"itemId"`
	Action   string              `json:"action"`
	Status   loanStatus          `json:"status"`
	LoanDate string              `json:"loanDate"`
	DueDate  string              `json:"dueDate"`
	Item     circulationLoanItem `json:"item"`
}

func countLoanStorageFiles(filepath string) (numMatching int) {
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		logger.Fatal(err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), "loan-storage-loans-") {
			numMatching++
			// fmt.Println(f.Name())
		}
	}
	return numMatching
}

// Return inventoryItems.json as a map, indexed by item ID
func makeItemsMap(filepath string) map[string]inventoryItem {
	itemsMap := make(map[string]inventoryItem)
	itemsChnl := make(chan interface{}, 1)
	go streamFolioSliceItem("items", filepath, itemsChnl)
	for oneItem := range itemsChnl {
		var itemObj inventoryItem
		mapstructure.Decode(oneItem, &itemObj)
		key := itemObj.ID
		itemsMap[key] = itemObj
	}
	return itemsMap
}

// GenerateCirculationLoans makes the same number of loans as found in loans.json
func GenerateCirculationLoans(filedef FileDef, outputParams OutputParams) {
	itemsPath := filepath.Join(outputParams.OutputDir, "inventory-items-1.json")
	itemsMap := makeItemsMap(itemsPath)
	numFiles := countLoanStorageFiles(outputParams.OutputDir)
	numThings := 0
	for i := 1; i <= numFiles; i++ {
		var circLoans []interface{}
		loanChnl := streamOutputLinearly(outputParams, "loan-storage-loans-"+strconv.Itoa(i)+".json", "loans")
		for oneLoan := range loanChnl {
			var loanObj circulationLoan
			mapstructure.Decode(oneLoan, &loanObj)
			itemID := loanObj.ItemID
			oneItem := itemsMap[itemID]
			loanObj.Item = circulationLoanItem{
				Title:            oneItem.Title,
				Barcode:          oneItem.Barcode,
				HoldingsRecordID: oneItem.HoldingsRecordID,
				Status:           oneItem.Status,
				Location: itemStatus{
					Name: oneItem.EffectiveLocation.Name,
				},
			}
			circLoans = append(circLoans, loanObj)
		}
		writeOutput(outputParams, fileNumStr(filedef, i), filedef.ObjectKey, circLoans)
		numThings += len(circLoans)
	}
	filedef.NumFiles = numFiles
	filedef.N = numThings
	updateManifest(filedef, outputParams)
}
