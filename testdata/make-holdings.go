package testdata

import (
	uuid "github.com/satori/go.uuid"
)

type holding struct {
	ID string `json:"id"`
}

func GenerateHoldings(allParams AllParams, numHoldings int) {
	makeHolding := func() materialType {
		return materialType{
			ID: uuid.Must(uuid.NewV4()).String(),
		}
	}
	var holdings []interface{}
	for i := 0; i < numHoldings; i++ {
		h := makeHolding()
		holdings = append(holdings, h)
	}

	filename := "holdings.json"
	objKey := "holdingsRecords"
	writeOutput(allParams.Output, filename, objKey, holdings)

	updateManifest(FileDef{
		Module:    "mod-inventory-storage",
		Path:      "/holdings-storage/holdings",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  1,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/holdings-storage.html",
		N:         numHoldings,
	}, allParams.Output)
}
