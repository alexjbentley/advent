package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"utils"
)

func readDepths() []int {
	data, err := os.ReadFile("files/depths.txt")
	utils.Check(err)

	readingStrings := strings.Split(string(data), "\n")
	numReadings := len(readingStrings)
	readings := make([]int, numReadings)
	for i, readingString := range readingStrings {
		readings[i], err = strconv.Atoi(readingString)
		utils.Check(err)
	}
	return readings

}

func sumWindow(window []int) int {
	return window[0] + window[1] + window[2]
}

func main() {
	readings := readDepths()
	numIncreased := 0
	leadWindow := make([]int, 0) // Window at the front
	prevWindow := make([]int, 0) // Trailing window
	prevReading := 0
	for i, reading := range readings {
		if i == 0 {
			// First load, don't alter the trailing window
			leadWindow = append(leadWindow, reading) // Push to lead
		} else if i < 3 {
			// We haven't loaded enough values to compare yet
			leadWindow = append(leadWindow, reading)     // Push to lead
			prevWindow = append(prevWindow, prevReading) // Push to trailing
		} else if i == 3 {
			// Special case, leadWindow will end up with 4 elements before the pop
			leadWindow = append(leadWindow, reading)     // Push to lead
			prevWindow = append(prevWindow, prevReading) // Push to trailing
			leadWindow = leadWindow[1:]
			if sumWindow(leadWindow) > sumWindow(prevWindow) {
				numIncreased++
			}
		} else {
			// Sliding windows are ready
			// Push
			leadWindow = append(leadWindow, reading)     // Push to lead
			prevWindow = append(prevWindow, prevReading) // Push to trailing
			leadWindow = leadWindow[1:]
			prevWindow = prevWindow[1:]

			// Compare
			if sumWindow(leadWindow) > sumWindow(prevWindow) {
				numIncreased++
			}
		}
		prevReading = reading
	}
	fmt.Println(numIncreased)
}
