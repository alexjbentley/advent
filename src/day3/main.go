package main

import (
	"fmt"
	"os"
	"strings"
	"utils"
)

func getMask(length int) (mask int) {
	mask = 0b1
	for i := 0; i < length; i++ {
		mask = mask | 0b1<<i
	}
	return
}

func readReport(filename string) (bitCounts []int, reportLength int) {
	data, err := os.ReadFile(filename)
	utils.Check(err)

	// data is currently a []byte
	reportLines := strings.Split(string(data), "\n")
	reportLength = len(reportLines)

	// Let's get the length of the lines
	numBits := len(reportLines[0])

	// Intialise our slice to be zeroed
	bitCounts = make([]int, numBits)
	for i := 0; i < numBits; i++ {
		bitCounts[i] = 0
	}

	// For each report line
	for _, reportLine := range reportLines {
		// Split into characters
		lineRunes := strings.Split(reportLine, "")
		// Note that MSB will be read first!
		for i := 0; i < len(lineRunes); i++ {
			if lineRunes[i] == "1" {
				bitCounts[i]++
			}
		}
	}

	return
}

func main() {
	bitCounts, reportLength := readReport(os.Args[1])
	gammaRate := 0b0
	numBits := len(bitCounts)

	// Note that bitCounts[0] concerns MSB!
	for i := 0; i < numBits; i++ {
		if bitCounts[i] > reportLength/2 {
			gammaRate = gammaRate | 0b1<<(numBits-i-1)
		}
	}
	epsilonRate := gammaRate ^ getMask(numBits)
	fmt.Printf("Gamma: %012b (%d)\n", gammaRate, gammaRate)
	fmt.Printf("Epsilon: %012b (%d)\n", epsilonRate, epsilonRate)

	powerConsumption := gammaRate * epsilonRate
	fmt.Println(powerConsumption)
}
