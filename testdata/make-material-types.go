package testdata

import (
	uuid "github.com/satori/go.uuid"
)

type materialType struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Source string `json:"source"`
}

func GenerateMaterialTypes(filedef FileDef, outputParams OutputParams) {
	typeList := []string{"dvd", "video recording", "microform", "electronic resource", "text",
		"sound recording", "unspecified", "book"}
	makeMaterialType := func(typeName string) materialType {
		return materialType{
			Name:   typeName,
			ID:     uuid.Must(uuid.NewV4()).String(),
			Source: "folio",
		}
	}
	var types []interface{}
	for i := 0; i < len(typeList); i++ {
		l := makeMaterialType(typeList[i])
		types = append(types, l)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, types)
	filedef.N = len(typeList)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
