package testdata

import (
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), "loans.json.") {
			numMatching++
			// fmt.Println(f.Name())
		}
	}
	return numMatching
}

func makeItemsMap(filepath string) map[string]item {
	itemsMap := make(map[string]item)
	itemsChnl := make(chan interface{}, 1)
	go streamSliceItem(filepath, itemsChnl)
	for oneItem := range itemsChnl {
		var itemObj item
		mapstructure.Decode(oneItem, &itemObj)
		key := itemObj.ID
		itemsMap[key] = itemObj
	}
	return itemsMap
}

func GenerateCirculationLoans(outputDir string) {
	var circLoans []interface{}
	itemsPath := filepath.Join(outputDir, "items.json")
	itemsMap := makeItemsMap(itemsPath)
	numFiles := countLoanStorageFiles(outputDir)
	for i := 1; i <= numFiles; i++ {
		loanStorageFilepath := filepath.Join(outputDir, "loans.json."+strconv.Itoa(i))
		loanChnl := make(chan interface{}, 1)
		go streamSliceItem(loanStorageFilepath, loanChnl)
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
		filepath := filepath.Join(outputDir, "circloan.json."+strconv.Itoa(i))
		fmt.Println("Writing", filepath)
		writeSliceToFile(filepath, circLoans, true)
	}
}