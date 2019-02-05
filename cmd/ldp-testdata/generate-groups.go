package main

import (
	uuid "github.com/satori/go.uuid"
)

type groupMetadata struct {
	CreatedDate     string `json:"createdDate"`
	CreatedByUserID string `json:"createdByUserId"`
	UpdatedDate     string `json:"updatedDate"`
	UpdatedByUserID string `json:"updatedByUserId"`
}
type group struct {
	Group    string        `json:"group"`
	Desc     string        `json:"desc"`
	ID       string        `json:"id"`
	Metadata groupMetadata `json:"metadata"`
}

func generateGroups(filepath string) {
	groupNames := []string{
		"Freshman", "Sophomore", "Junior", "Senior", "Graduate", "Alumni", "Faculty",
		"Staff", "Affiliate_A", "Affiliate_B", "Affiliate_C", "Affiliate_D",
	}
	groupDesc := []string{
		"First year in undergrad", "Second year in undergrad", "Third year in undergrad",
		"Fourth year in undergrad", "Masters or Doctoral student", "Graduated from undergrad",
		"Professor or other faculty", "Staff at the university", "Type of affiliate",
		"Type of affiliate", "Type of affiliate", "Type of affiliate",
	}
	randomDate := "2018-11-19T14:29:58.542+0000"

	makeGroup := func(i int) group {
		creator := uuid.Must(uuid.NewV4()).String()
		return group{
			Group: groupNames[i],
			Desc:  groupDesc[i],
			ID:    uuid.Must(uuid.NewV4()).String(),
			Metadata: groupMetadata{
				CreatedDate:     randomDate,
				CreatedByUserID: creator,
				UpdatedDate:     randomDate,
				UpdatedByUserID: creator,
			}}
	}

	var groups []interface{}
	for i := 0; i < len(groupNames); i++ {
		g := makeGroup(i)
		groups = append(groups, g)
	}
	writeSliceToFile(filepath, groups, true)
}
