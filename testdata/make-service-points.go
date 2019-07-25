package testdata

type holdShelfExpiryPeriod struct {
	Duration   int    `json:"duration"`
	IntervalID string `json:"intervalId"`
}

type servicePoint struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	Code                  string                `json:"code"`
	DiscoveryDisplayName  string                `json:"discoveryDisplayName"`
	ShelvingLagTime       int                   `json:"shelvingLagTime"`
	PickupLocation        bool                  `json:"pickupLocation"`
	HoldShelfExpiryPeriod holdShelfExpiryPeriod `json:"holdShelfExpiryPeriod"`
	StaffSlips            []string              `json:"staffSlips"`
}

func GenerateServicePoints(filedef FileDef, outputParams OutputParams) {
	servicePointLiterals := []servicePoint{
		servicePoint{
			ID:                   "7c5abc9f-f3d7-4856-b8d7-6712462ca007",
			Name:                 "Online",
			Code:                 "Online",
			DiscoveryDisplayName: "Online",
			ShelvingLagTime:      0,
			PickupLocation:       false,
			StaffSlips:           []string{},
		},
		servicePoint{
			ID:                   "c4c90014-c8c9-4ade-8f24-b5e313319f4b",
			Name:                 "Circ Desk 2",
			Code:                 "cd2",
			DiscoveryDisplayName: "Circulation Desk -- Back Entrance",
			PickupLocation:       true,
			HoldShelfExpiryPeriod: holdShelfExpiryPeriod{
				Duration:   5,
				IntervalID: "Days",
			},
			StaffSlips: []string{},
		},
		servicePoint{
			ID:                   "3a40852d-49fd-4df2-a1f9-6e2641a6e91f",
			Name:                 "Circ Desk 1",
			Code:                 "cd1",
			DiscoveryDisplayName: "Circulation Desk -- Hallway",
			PickupLocation:       true,
			HoldShelfExpiryPeriod: holdShelfExpiryPeriod{
				Duration:   3,
				IntervalID: "Weeks",
			},
			StaffSlips: []string{},
		},
	}

	var servicePoints []interface{}
	for i := 0; i < len(servicePointLiterals); i++ {
		servicePoints = append(servicePoints, servicePointLiterals[i])
	}

	writeOutput(outputParams, fileNumStr(filedef, 1), filedef.ObjectKey, servicePoints)
	filedef.NumFiles = 1
	updateManifest(filedef, outputParams)
}
