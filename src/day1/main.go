package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"utils"
)

func main() {
	data, err := os.ReadFile("files/depths.txt")
	utils.Check(err)

	readingStrings := strings.Split(string(data), "\n")
	numReadings := len(readingStrings)
	readings := make([]int, numReadings)
	for i, readingString := range readingStrings {
		readings[i], err = strconv.Atoi(readingString)
		utils.Check(err)
	}
	prevReading := 0
	numIncreased := -1 // First will be a false positive

	for _, reading := range readings {
		if reading > prevReading {
			numIncreased++
		}
		prevReading = reading
	}
	fmt.Println(numIncreased)
}
