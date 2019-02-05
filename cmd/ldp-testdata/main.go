package main

import (
	"fmt"
)

func main() {
	fmt.Println("ldp-testdata")

	// generateGroups("./groups.json")
	// generateUsers("./users.json")
	// generateLocations("./locations.json")
	// generateItems("./items.json")
	tenMillion := 10000000
	generateLoans("./loans.json", tenMillion)
	// groupsFilename := "extract-output/sample/groups.json"
	// groupsChnl := make(chan string)
	// go streamRandomLine(groupsFilename, groupsChnl)
	// newVal, ok := <-groupsChnl
	// newVal2, _ := <-groupsChnl
	// if ok {
	// 	fmt.Println(newVal, newVal2)
	// }
	// threeGenerators()
}
