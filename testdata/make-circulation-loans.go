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
		if strings.HasPrefix(f.Name(), "loans.json.") {
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
func GenerateCirculationLoans(allParams AllParams, ignore int) {
	outputParams := allParams.Output
	var circLoans []interface{}
	filename := "circloan.json"
	objKey := "loans"
	itemsPath := filepath.Join(outputParams.OutputDir, "inventoryItems.json")
	itemsMap := makeItemsMap(itemsPath)
	numFiles := countLoanStorageFiles(outputParams.OutputDir)
	numThings := 0
	for i := 1; i <= numFiles; i++ {
		loanChnl := streamOutputLinearly(outputParams, "loans.json."+strconv.Itoa(i), "loans")
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
		writeOutput(outputParams, filename+"."+strconv.Itoa(i), objKey, circLoans)
		numThings += len(circLoans)
	}
	updateManifest(FileDef{
		Module:    "mod-circulation",
		Path:      "/circulation/loans",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  numFiles,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-circulation/circulation.html",
		N:         numThings,
	}, allParams.Output)
}
