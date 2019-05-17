package testdata

import (
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
)

type user struct {
	Username    string   `json:"username"`
	ID          string   `json:"id"`
	Barcode     string   `json:"barcode"`
	Active      bool     `json:"active"`
	Type        string   `json:"type"`
	PatronGroup string   `json:"patronGroup"`
	ProxyFor    []string `json:"proxyFor"`
}

func isActive() bool {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(2)
	if randNum == 0 {
		return true
	}
	return false
}

func GenerateUsers(filedef FileDef, outputParams OutputParams) {
	chnl := streamRandomItem(outputParams, "groups-1.json", "usergroups")
	makeUser := func() user {
		randomGroup, _ := <-chnl
		randomGroupMap := randomGroup.(map[string]interface{})
		randomGroupID := randomGroupMap["id"].(string)
		return user{
			Username:    fake.UserName(),
			ID:          uuid.Must(uuid.NewV4()).String(),
			Barcode:     fake.DigitsN(16),
			Active:      isActive(),
			Type:        "patron",
			PatronGroup: randomGroupID,
			ProxyFor:    make([]string, 0),
		}
	}
	// fmt.Printf("%+v\n", makeUser())
	var users []interface{}
	for i := 0; i < filedef.N; i++ {
		u := makeUser()
		users = append(users, u)
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, users)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
