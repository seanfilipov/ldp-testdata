package testdata

import (
	"path"
	"runtime"

	"github.com/icrowley/fake"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
)

// RAML: https://github.com/folio-org/mod-inventory-storage/blob/master/ramls/instance.json

type contributer struct {
	Name                  string `json:"name"`
	ContributorNameTypeID string `json:"contributorNameTypeId"`
	Primary               bool   `json:"primary"`
}
type instance struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Source         string        `json:"source"`
	Contributors   []contributer `json:"contributors"`
	InstanceTypeID string        `json:"instanceTypeId"`
}

func GenerateInstances(filedef FileDef, outputParams OutputParams) {

	typeChnl := streamRandomItem(outputParams, "instance-types-1.json", "instanceTypes")
	bookChnl := make(chan string, 1)
	_, nameOfThisFile, _, _ := runtime.Caller(0)
	pkgDir := path.Dir(nameOfThisFile)
	go streamRandomLine(pkgDir+"/book_titles.txt", bookChnl)

	makeInstance := func() instance {
		randomBookTitle, _ := <-bookChnl
		randomType, _ := <-typeChnl
		var instanceTypeObj instanceType
		mapstructure.Decode(randomType, &instanceTypeObj)

		return instance{
			ID:     uuid.Must(uuid.NewV4()).String(),
			Title:  randomBookTitle,
			Source: "MARC",
			Contributors: []contributer{
				contributer{
					Name: fake.FullName(),
				},
			},
			InstanceTypeID: instanceTypeObj.ID,
		}
	}
	var instances []interface{}
	for i := 0; i < filedef.N; i++ {
		oneInstance := makeInstance()
		instances = append(instances, oneInstance)
	}
	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, instances)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
