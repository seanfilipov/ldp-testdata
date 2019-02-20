package testdata

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// generate random values based on either a list of values or a range/step
//   - A list of values would be like the values in groups.groupname:
//       [faculty, staff, undergraduate student, graduate student]
//   - A range/step would something like a for-loop: range from 10-100, step 5 = 10, 15, 20, etc.

// outputSize corresponds to the number of rows in the table that you want to generate

// makeRandomValues doesn't know anything about the data.
// It just wants to know how much data to generate and how many possible values there are.

// Streams integers, each integer being between 0 and the given max
func makeRandomValues(maxValue, outputSize int, chnl chan int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < outputSize; i++ {
		chnl <- rand.Intn(maxValue)
	}
	close(chnl)
}

type ValueParams struct {
	MaxValue   int
	OutputSize int
	Channel    chan int
}
type ChannelResponse struct {
	Value int
	Ok    bool
}

// Example of streaming 3 random values and outputting a string with them
func threeInstances() {
	valueType1 := ValueParams{2, 10, make(chan int)}
	valueType2 := ValueParams{4, 10, make(chan int)}
	valueType3 := ValueParams{6, 10, make(chan int)}
	allValueTypes := []ValueParams{valueType1, valueType2, valueType3}
	for _, valueType := range allValueTypes {
		go makeRandomValues(valueType.MaxValue, valueType.OutputSize, valueType.Channel)
	}
	for {
		// 1) Get a response from each generator
		var responses []ChannelResponse
		for _, valueType := range allValueTypes {
			resp := ChannelResponse{}
			newVal, ok := <-valueType.Channel
			resp.Value = newVal
			resp.Ok = ok
			// In this example, the streamed response is added to an array
			// In actual usage, we'd want to use the response (e.g. forming a JSON record)
			responses = append(responses, resp)
		}
		// 2) Form a string from the responses
		var values []int
		var channelClosed = false
		for _, resp := range responses {
			if resp.Ok {
				values = append(values, resp.Value)
			} else {
				channelClosed = true
				break
			}
		}
		if channelClosed {
			break
		} else {
			delim := " "
			stringArray := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(values)), delim), "[]")
			fmt.Println("Received", stringArray)
		}
	}
}
