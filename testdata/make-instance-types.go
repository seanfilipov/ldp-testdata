package testdata

type instanceType struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Source string `json:"primary"`
}

func GenerateInstanceTypes(filedef FileDef, outputParams OutputParams) {
	typeLiterals := []instanceType{
		instanceType{
			ID:     "efe2e89b-0525-4535-aa9b-3ff1a131189e",
			Name:   "tactile image",
			Code:   "tci",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "225faa14-f9bf-4ecd-990d-69433c912434",
			Name:   "two-dimensional moving image",
			Code:   "tdi",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "a2c91e87-6bab-44d6-8adb-1fd02481fc4f",
			Name:   "other",
			Code:   "xxx",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "3e3039b7-fda0-4ac4-885a-022d457cb99c",
			Name:   "three-dimensional moving image",
			Code:   "tdm",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "3be24c14-3551-4180-9292-26a786649c8b",
			Name:   "performed music",
			Code:   "prm",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "e6a278fb-565a-4296-a7c5-8eb63d259522",
			Name:   "tactile notated movement",
			Code:   "tcn",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "6312d172-f0cf-40f6-b27d-9fa8feaf332f",
			Name:   "text",
			Code:   "txt",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "2022aa2e-bdde-4dc4-90bc-115e8894b8b3",
			Name:   "cartographic three-dimensional form",
			Code:   "crf",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "3363cdb1-e644-446c-82a4-dc3a1d4395b9",
			Name:   "cartographic dataset",
			Code:   "crd",
			Source: "rdacontent",
		},
		instanceType{
			ID:     "c208544b-9e28-44fa-a13c-f4093d72f798",
			Name:   "computer program",
			Code:   "cop",
			Source: "rdacontent",
		},
	}

	var instanceTypes []interface{}
	for i := 0; i < len(typeLiterals); i++ {
		instanceTypes = append(instanceTypes, typeLiterals[i])
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, instanceTypes)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
