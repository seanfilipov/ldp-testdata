package testdata

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// The purpose of this file is to provide the randomCallNumber() function

// 21 very broad categories defined by the Library of Congress classification system
var firstLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "Z"}

// all 26 letters
var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func getFirstLetter() string {
	index := random(0, len(firstLetters)-1)
	return firstLetters[index]
}
func getLetter() string {
	index := random(0, len(alphabet)-1)
	return alphabet[index]
}
func getNumbers() string {
	var str strings.Builder
	for i := 0; i < random(2, 4); i++ {
		str.WriteString(strconv.Itoa(random(0, 9)))
	}
	return str.String()
}

// examples: ".C22", ".A1", ".Z378"
func getDecimalPart() string {
	var str strings.Builder
	str.WriteString("." + getLetter())
	for i := 0; i < random(1, 3); i++ {
		str.WriteString(strconv.Itoa(random(0, 9)))
	}
	return str.String()
}

func getYear() string {
	return strconv.Itoa(random(1920, 2019))
}

func randomCallNumber() string {
	rand.Seed(time.Now().UnixNano())
	firstTwoLetters := getFirstLetter() + getLetter()
	middleNumbers := getNumbers()
	return firstTwoLetters + middleNumbers + " " + getDecimalPart() + " " + getYear()
}
