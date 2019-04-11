package testdata

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/folio-org/ldp-testdata/logging"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

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
	outputParams OutputParams
	Filename     string
	ObjectKey    string
	ItemChnl     chan interface{}
	UserChnl     chan interface{}
	CheckedOut   map[string]loan
	EndDay       int
	TxnPerDay    int
	TxnPerFile   int
}

// Split loans into separate files
// Spread out

// Get a random item ID
// If that item has already been checked out, check it back in
// Otherwise, check it out
func (lg loanGenerator) makeLoanTxn(date time.Time, itemID string) (retLoan loan) {
	if itemID == "" {
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
	}
	if checkedOutLoan, ok := lg.CheckedOut[itemID]; ok {
		retLoan = checkedOutLoan
		retLoan.ID = uuid.Must(uuid.NewV4()).String()
		retLoan.Action = "checkedin"
		retLoan.Status.Name = "Closed"
		delete(lg.CheckedOut, itemID)
	} else {
		randomUser, _ := <-lg.UserChnl
		var userObj user
		mapstructure.Decode(randomUser, &userObj)
		l := loan{
			ID:       uuid.Must(uuid.NewV4()).String(),
			UserID:   userObj.ID,
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

// Write until the number of txnsPerFile is reached
// OR the number of days is reached
func (lg loanGenerator) makeLoans(startDay int) (day int, loans []interface{}) {
	layout := time.RFC3339
	jan1 := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Loop over the number of days
	for day = startDay; day < lg.EndDay; day++ {
		numCheckins := 0
		date := jan1.Add(time.Hour * 24 * time.Duration(day)) // transform it into a date format
		for itemID, loanObj := range lg.CheckedOut {
			dueDate, _ := time.Parse(layout, loanObj.DueDate)
			if dueDate == date {
				loans = append(loans, lg.makeLoanTxn(date, itemID))
				numCheckins++
			}
		}
		numCheckouts := lg.TxnPerDay - numCheckins
		// fmt.Printf("day:%d numCheckouts: %d, TxnPerDay: %d, numCheckins: %d\n", day, numCheckouts, lg.TxnPerDay, numCheckins)
		// Loop over the number of txnsPerFile
		for i := 0; i < numCheckouts; i++ {
			loans = append(loans, lg.makeLoanTxn(date, ""))
			if len(loans) >= lg.TxnPerFile {
				return
			}
		}
	}
	return lg.EndDay, loans
}

func (lg loanGenerator) generateLoansSingleFile(startDay, callNum int) (int, int) {
	reachedDay, loans := lg.makeLoans(startDay)
	callNumStr := strconv.Itoa(callNum)
	filename := lg.Filename + "." + callNumStr
	writeOutput(lg.outputParams, filename, lg.ObjectKey, loans)
	totalWritten := strconv.Itoa(((callNum - 1) * lg.TxnPerFile) + len(loans))
	logger.Debugf("Wrote %d transactions to %s (%s total)\n", len(loans), filename, totalWritten)
	return reachedDay, len(loans)
}
func (lg loanGenerator) run() (counter, totalLoansMade int) {
	runCount := 0
	reachedDay := 0
	for reachedDay != lg.EndDay {
		var numLoans int
		runCount++
		reachedDay, numLoans = lg.generateLoansSingleFile(reachedDay, runCount)
		logger.Debugf("Run #%d: reached day %d\n", runCount, reachedDay)
		counter++
		totalLoansMade += numLoans
	}
	return
}

func GenerateLoans(allParams AllParams, totalNumTxns int) {
	outputParams := allParams.Output
	numDays := 365
	txnPerFile := 100000
	txnPerDay := int(math.Ceil(float64(totalNumTxns / numDays)))

	numFilesNeeded := strconv.Itoa(int(math.Ceil(float64((txnPerDay * numDays) / txnPerFile))))
	logger.Debug("Going to write ~" + numFilesNeeded + " files")
	filename := "loans.json"
	objKey := "loans"
	lg := loanGenerator{
		outputParams,
		filename,
		objKey,
		streamRandomItem(outputParams, "storageItems.json", "items"),
		streamRandomItem(outputParams, "users.json", "users"),
		make(map[string]loan),
		numDays,
		txnPerDay,
		txnPerFile,
	}
	numFiles, totalLoansMade := lg.run()
	updateManifest(FileDef{
		Module:    "mod-circulation-storage",
		Path:      "/loan-storage/loans",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  numFiles,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-circulation-storage/loan-storage.html",
		N:         totalLoansMade,
	}, outputParams)
}
