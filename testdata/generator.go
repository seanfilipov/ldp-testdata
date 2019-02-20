package testdata

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// generate sample data that are random, based on either a list of values or a range/step
//   - A list of values would be like the values in groups.groupname:
//       [faculty, staff, undergraduate student, graduate student]
//   - A range/step would something like a for-loop: range from 10-100, step 5 = 10, 15, 20, etc.

// outputSize corresponds to the number of rows in the table that you want to generate

// generateRandom doesn't know anything about the data.
// It just wants to know how much data to generate and how many possible values there are.

// Streams integers, each integer being between 0 and the given max
func GenerateRandomValues(maxValue, outputSize int, chnl chan int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < outputSize; i++ {
		chnl <- rand.Intn(maxValue)
	}
	close(chnl)
}

type Generator struct {
	MaxValue   int
	OutputSize int
	Channel    chan int
}
type ChannelResponse struct {
	Value int
	Ok    bool
}

// Creates 3 example generators, streams responses from each
func threeGenerators() {
	firstGen := Generator{2, 10, make(chan int)}
	secondGen := Generator{4, 10, make(chan int)}
	thirdGen := Generator{6, 10, make(chan int)}
	generators := []Generator{firstGen, secondGen, thirdGen}
	for _, gen := range generators {
		go GenerateRandomValues(gen.MaxValue, gen.OutputSize, gen.Channel)
	}
	for {
		// Get a response from each generator
		var responses []ChannelResponse
		for _, gen := range generators {
			resp := ChannelResponse{}
			newVal, ok := <-gen.Channel
			resp.Value = newVal
			resp.Ok = ok
			// In this example, the streamed response is added to an array
			// In actual usage, we'd want to use the response (e.g. forming a JSON record)
			responses = append(responses, resp)
		}
		// Form a string from the responses
		var genValues []int
		var channelClosed = false
		for _, resp := range responses {
			if resp.Ok {
				genValues = append(genValues, resp.Value)
			} else {
				channelClosed = true
				break
			}
		}
		if channelClosed {
			break
		} else {
			delim := " "
			stringArray := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(genValues)), delim), "[]")
			fmt.Println("Received", stringArray)
		}
	}
}
