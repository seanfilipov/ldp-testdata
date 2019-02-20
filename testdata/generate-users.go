package testdata

import (
	"math/rand"
	"path/filepath"
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
	} else {
		return false
	}
}

func GenerateUsers(outputDir string, numUsers int) {
	chnl := make(chan interface{}, 1)
	groupsPath := filepath.Join(outputDir, "groups.json")
	go streamRandomSliceItem(groupsPath, chnl)
	makeUser := func() user {
		randomGroup, _ := <-chnl
		randomGroupMap := randomGroup.(map[string]interface{})
		randomGroupName := randomGroupMap["group"].(string)
		return user{
			Username:    fake.UserName(),
			ID:          uuid.Must(uuid.NewV4()).String(),
			Barcode:     fake.DigitsN(16),
			Active:      isActive(),
			Type:        "patron",
			PatronGroup: randomGroupName,
			ProxyFor:    make([]string, 0),
		}
	}
	// fmt.Printf("%+v\n", makeUser())
	var users []interface{}
	for i := 0; i < numUsers; i++ {
		u := makeUser()
		users = append(users, u)
	}
	filepath := filepath.Join(outputDir, "users.json")
	writeSliceToFile(filepath, users, true)
}
