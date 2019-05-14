package testdata

import (
	uuid "github.com/satori/go.uuid"
)

type materialType struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Source string `json:"source"`
}

func GenerateMaterialTypes(allParams AllParams, numMaterialTypes int) {
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

	filename := "material-types.json"
	objKey := "mtypes"
	writeOutput(allParams.Output, filename, objKey, types)

	updateManifest(FileDef{
		Module:    "mod-inventory-storage",
		Path:      "/material-types",
		Filename:  filename,
		ObjectKey: objKey,
		NumFiles:  1,
		Doc:       "https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/material-type.html",
		N:         numMaterialTypes,
	}, allParams.Output)
}
