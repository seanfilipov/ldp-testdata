package testdata

import "encoding/json"

type addressType struct {
	ID          string `json:"id"`
	AddressType string `json:"addressType"`
	Desc        string `json:"desc"`
}

type addressTypesFile struct {
	AddressTypes []addressType `json:"addressTypes"`
}

func GenerateAddressTypes(filedef FileDef, outputParams OutputParams) {
	addressTypeLiterals := []addressType{
		addressType{
			ID:          "93d3d88d-499b-45d0-9bc7-ac73c3a19880",
			AddressType: "Home",
			Desc:        "Home Address",
		},
		addressType{
			ID:          "1c4b225f-f669-4e9b-afcd-ebc0e273a34e",
			AddressType: "Work",
			Desc:        "Work Address",
		},
	}

	var addressTypes []interface{}
	for i := 0; i < len(addressTypeLiterals); i++ {
		addressTypes = append(addressTypes, addressTypeLiterals[i])
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, addressTypes)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}

// usage: addressType := readAddressTypes(outputParams, "addresstypes-1.json")
func readAddressTypes(params OutputParams, filename string) []addressType {
	byteValue := readFile(params, filename)
	var addressTypeFileObj addressTypesFile
	json.Unmarshal(byteValue, &addressTypeFileObj)
	return addressTypeFileObj.AddressTypes
}
