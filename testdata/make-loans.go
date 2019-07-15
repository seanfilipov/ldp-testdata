package testdata

import (
	"math"
	"os"
	"time"

	"github.com/folio-org/ldp-testdata/logging"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

// https://github.com/folio-org/mod-circulation-storage/blob/master/ramls/loan.json

var logger = logging.Logger

type loanStatus struct {
	Name string `json:"name"`
}
type loan struct {
	ID       string     `json:"id"`
	UserID   string     `json:"userId"`
	ItemID   string     `json:"itemId"`
	Action   string     `json:"action"`
	Status   loanStatus `json:"status"`
	LoanDate string     `json:"loanDate"`
	DueDate  string     `json:"dueDate"`
}

type loanGenerator struct {
	ItemChnl   chan interface{}
	UserChnl   chan interface{}
	CheckedOut map[string]loan
}

func (lg loanGenerator) randomItemID() (itemID string) {
	randomItem, ok := <-lg.ItemChnl
	if !ok {
		logger.Error("Could not get item from channel")
	}
	var itemObj storageItem
	mapstructure.Decode(randomItem, &itemObj)
	itemID = itemObj.ID
	if itemID == "" {
		close(lg.ItemChnl)
		logger.Errorf("Item received from channel has no ID field: %s", randomItem)
		os.Exit(1)
	}
	return
}
func (lg loanGenerator) randomUserID() (userID string) {
	randomUser, _ := <-lg.UserChnl
	var userObj user
	mapstructure.Decode(randomUser, &userObj)
	return userObj.ID
}

// 1. Get a random item ID
// 2. If that item has already been checked out, check it back in
// 3. Otherwise, check it out
func (lg loanGenerator) makeLoanTxn(date time.Time) (retLoan loan) {
	itemID := lg.randomItemID()
	if checkedOutLoan, ok := lg.CheckedOut[itemID]; ok {
		retLoan = checkedOutLoan
		retLoan.ID = uuid.Must(uuid.NewV4()).String()
		retLoan.Action = "checkedin"
		retLoan.Status.Name = "Closed"
		delete(lg.CheckedOut, itemID)
	} else {
		l := loan{
			ID:       uuid.Must(uuid.NewV4()).String(),
			UserID:   lg.randomUserID(),
			ItemID:   itemID,
			Action:   "checkedout",
			Status:   loanStatus{Name: "Open"},
			LoanDate: date.Format(time.RFC3339),
			DueDate:  date.Add(time.Hour * 24 * 7 * 2).Format(time.RFC3339), // loan duration: 14 days
		}
		lg.CheckedOut[itemID] = l
		retLoan = l
	}
	return retLoan
}

func GenerateLoans(filedef FileDef, outputParams OutputParams) {
	lg := loanGenerator{
		streamRandomItem(outputParams, "item-storage-items-1.json", "items"),
		streamRandomItem(outputParams, "users-1.json", "users"),
		make(map[string]loan),
	}

	N := filedef.N
	numFilesWritten := 0
	numDays := 365 // approximate; because N cannot be evenly divided into 365 days, the remainder goes into overflow days
	nInFile := 0
	nInDay := 0
	maxNInFile := 100000
	maxNInDay := int(math.Ceil(float64(N / numDays)))
	jan1 := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	today := 0                                                   // dayNum (0..365)
	todayDate := jan1.Add(time.Hour * 24 * time.Duration(today)) // transform it into a date format

	loans := make([]interface{}, 0)
	for i := 0; i < N; i++ {
		loans = append(loans, lg.makeLoanTxn(todayDate))
		nInDay++
		nInFile++
		logger.Debugln("nInFile =", nInFile, todayDate)
		if nInFile == maxNInFile || nInFile == N {
			numFilesWritten++
			filename := fileNumStr(filedef, numFilesWritten)
			writeOutput(outputParams, filename, "loans", loans)
			loans = make([]interface{}, 0)
			nInFile = 0
		}
		if nInDay == maxNInDay {
			today++
			todayDate = jan1.Add(time.Hour * 24 * time.Duration(today)) // transform it into a date format
			nInDay = 0
		}
	}

	filedef.NumFiles = numFilesWritten
	updateManifest(filedef, outputParams)
}
