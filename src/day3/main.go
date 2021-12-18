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

func readReport() (bitCounts []int, reportLength int) {
	data, err := os.ReadFile("files/diagnostic-report.txt")
	utils.Check(err)

	// data is currently a []byte
	reportLines := strings.Split(string(data), "\n")
	reportLength = len(reportLines)

	// Let's get the length of the lines
	numBits := len(reportLines[0])
	bitCounts = make([]int, numBits)
	for i := 0; i < numBits; i++ {
		bitCounts[i] = 0
	}

	for _, reportLine := range reportLines {
		lineRunes := strings.Split(reportLine, "")
		for i := 0; i < numBits; i++ {
			if lineRunes[i] == "1" {
				bitCounts[i]++
			}
		}
	}

	return
}

func main() {
	bitCounts, reportLength := readReport()
	gammaRate := 0b0
	for i := 0; i < len(bitCounts); i++ {
		if bitCounts[i] > reportLength/2 {
			gammaRate = gammaRate | 0b1<<i
		}
	}
	epsilonRate := gammaRate ^ getMask(len(bitCounts))
	fmt.Printf("Gamma: %012b %d\n", gammaRate, gammaRate)
	fmt.Printf("Epsilon: %012b %d\n", epsilonRate, epsilonRate)

	powerConsumption := gammaRate * epsilonRate
	fmt.Println(powerConsumption)
}
