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
	Personal    personal `json:"personal"`
}
type personal struct {
	Lastname    string    `json:"lastName"`
	Firstname   string    `json:"firstName"`
	Middlename  string    `json:"middleName"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	MobilePhone string    `json:"mobilePhone"`
	DateOfBirth string    `json:"dateOfBirth"`
	Addresses   []address `json:"addresses"`
}
type address struct {
	ID             string `json:"id"`
	CountryID      string `json:"countryId"`
	AddressLine1   string `json:"addressLine1"`
	AddressLine2   string `json:"addressLine2"`
	City           string `json:"city"`
	Region         string `json:"region"`
	PostalCode     string `json:"postalCode"`
	AddressTypeID  string `json:"addressTypeId"`
	PrimaryAddress bool   `json:"primaryAddress"`
}

func isActive() bool {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(2)
	if randNum == 0 {
		return true
	}
	return false
}
func randomDate() string {
	year := random(1950, 2010)
	month := time.Month(random(1, 12))
	day := random(0, 28)
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return date.Format(time.RFC3339)
}
func makePersonal(addrTypeID string) personal {
	return personal{
		Lastname:    fake.LastName(),
		Firstname:   fake.FirstName(),
		Middlename:  fake.FirstName(),
		Email:       fake.EmailAddress(),
		Phone:       fake.Phone(),
		MobilePhone: fake.Phone(),
		DateOfBirth: randomDate(),
		Addresses: []address{
			address{
				ID:             uuid.Must(uuid.NewV4()).String(),
				CountryID:      uuid.Must(uuid.NewV4()).String(),
				AddressLine1:   fake.StreetAddress(),
				City:           fake.City(),
				Region:         fake.State(),
				PostalCode:     fake.Zip(),
				AddressTypeID:  addrTypeID,
				PrimaryAddress: true,
			},
		},
	}
}

func GenerateUsers(filedef FileDef, outputParams OutputParams) {

	addressTypes := readAddressTypes(outputParams, "addresstypes-1.json")
	chnl := streamRandomItem(outputParams, "groups-1.json", "usergroups")
	makeUser := func() user {
		randomAddrType := addressTypes[random(0, len(addressTypes)-1)]
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
			Personal:    makePersonal(randomAddrType.ID),
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
