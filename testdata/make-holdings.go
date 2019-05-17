package testdata

import (
	uuid "github.com/satori/go.uuid"
)

type holding struct {
	ID string `json:"id"`
}

func GenerateHoldings(filedef FileDef, outputParams OutputParams) {
	makeHolding := func() materialType {
		return materialType{
			ID: uuid.Must(uuid.NewV4()).String(),
		}
	}
	var holdings []interface{}
	for i := 0; i < filedef.N; i++ {
		h := makeHolding()
		holdings = append(holdings, h)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, holdings)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
