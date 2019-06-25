package testdata

import (
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

// https://github.com/folio-org/mod-inventory-storage/blob/master/ramls/locationunit.raml
// https://github.com/folio-org/mod-inventory-storage/blob/master/ramls/examples/locinst.json

type locationUnit struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Code string `json:"code"`
}

func GenerateLocationUnitInstitutions(filedef FileDef, outputParams OutputParams) {
	makeLocationUnit := func() locationUnit {
		name := fake.LastName() + " Library"
		code := string(name[0]) + "L"
		return locationUnit{
			Name: name,
			ID:   uuid.Must(uuid.NewV4()).String(),
			Code: code,
		}
	}
	var locationUnits []interface{}
	for i := 0; i < filedef.N; i++ {
		l := makeLocationUnit()
		locationUnits = append(locationUnits, l)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, locationUnits)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
