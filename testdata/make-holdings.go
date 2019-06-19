package testdata

import (
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

type holding struct {
	ID                  string `json:"id"`
	InstanceID          string `json:"instanceId"`
	PermanentLocationID string `json:"permanentLocationId"`
}

func GenerateHoldings(filedef FileDef, outputParams OutputParams) {
	instanceChnl := streamOutputLinearly(outputParams, "instance-storage-instances-1.json", "instances")
	// numFiles := countFilesWithPrefix(outputParams.OutputDir, "instance-storage-instances")

	makeHolding := func(oneInstance interface{}) holding {
		var instanceObj instance
		mapstructure.Decode(oneInstance, &instanceObj)
		return holding{
			ID:         uuid.Must(uuid.NewV4()).String(),
			InstanceID: instanceObj.ID,
		}
	}
	var holdings []interface{}
	for oneInstance := range instanceChnl {
		h := makeHolding(oneInstance)
		holdings = append(holdings, h)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, holdings)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
