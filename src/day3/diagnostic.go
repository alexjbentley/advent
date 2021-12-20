package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"utils"
)

var ZERO = string("0")
var ONE = string("1")
var OXYGEN = true
var CO2 = false

type Node struct {
	// TODO Need to have our own bitCounts for the values we possess. Do away with numBits, we can get it from our own
	// bit counts once we have them.
	values []string

	// The particular "bit" this node is concerned with. Analagous with depth.
	index     int
	bitCounts []int
	left      *Node
	right     *Node
}

func (node Node) filterValues(key string) (filtered []string) {
	filtered = []string{}
	for _, value := range node.values {
		bit := strings.Split(value, "")[node.index]
		if bit == key {
			filtered = append(filtered, value)
		}
	}
	return
}

func (node *Node) buildLeftNode() {
	filteredValues := node.filterValues(ZERO)
	fmt.Printf("Creating left node with values:\n")
	fmt.Println(filteredValues)
	if len(filteredValues) > 0 {
		node.left = &Node{filteredValues, node.index + 1, node.bitCounts, nil, nil}
		if len(filteredValues) > 1 {
			node.left.buildAllBelow()
		}
	}
}

func (node *Node) buildRightNode() {
	filteredValues := node.filterValues(ONE)
	fmt.Printf("Creating right node with values:\n")
	fmt.Println(filteredValues)
	if len(filteredValues) > 0 {
		node.right = &Node{filteredValues, node.index + 1, node.bitCounts, nil, nil}
		if len(filteredValues) > 1 {
			node.right.buildAllBelow()
		}
	}
}

// Recursively build all nodes below this one, based on the values
func (node *Node) buildAllBelow() {
	// TODO we shouldn't be storing any values in the non-leaf nodes. Memory efficiency!
	if node.index == len(node.bitCounts) || len(node.values) == 1 {
		// We're the leaf of the tree, stop.
		return
	}
	// Make left Node (0):
	node.buildLeftNode()

	// Make right Node (1):
	node.buildRightNode()
}

func (node Node) search(mode bool) []string {
	// Escape clause first!
	fmt.Printf("OXYGEN: %t\nNode with following being searched:\n%v\n", mode, node.values)
	// if len(node.values) == 1 || node.index == node.numBits-1 {
	if len(node.values) == 1 {
		return node.values
	}

	positiveCount := 0
	for _, value := range node.values {
		valueTokens := strings.Split(value, "")
		if valueTokens[node.index] == ONE {
			positiveCount++
		}
	}
	if mode == OXYGEN {
		// Mode is OXYGEN, look for the most common, preferring 1.
		if float64(positiveCount) >= float64((len(node.values)))/2 {
			return node.right.search(mode)
		} else {
			return node.left.search(mode)
		}
	} else {
		// Mode is CO2, look for the least common, preferring 0.
		if float64(positiveCount) >= float64((len(node.values)))/2 {
			return node.left.search(mode)
		} else {
			return node.right.search(mode)
		}
	}
}

type DiagReport struct {
	gammaRate, epsilonRate int
	reportLines            []string
	rootNode               *Node
	bitCounts              []int
}

func (diagReport *DiagReport) initFromFile(filename string) {
	data, err := os.ReadFile(filename)
	utils.Check(err)
	diagReport.reportLines = strings.Split(string(data), "\n")
	numBits := len(diagReport.reportLines[0])

	// Intialise our slice to be zeroed
	diagReport.bitCounts = make([]int, numBits)
	for i := 0; i < numBits; i++ {
		diagReport.bitCounts[i] = 0
	}

	// For each report line
	for _, reportLine := range diagReport.reportLines {

		// Split into characters
		lineRunes := strings.Split(reportLine, "")
		// Note that MSB will be read first!
		for i := 0; i < len(lineRunes); i++ {
			if lineRunes[i] == ONE {
				diagReport.bitCounts[i]++
			}
		}
	}

	diagReport.rootNode = &Node{diagReport.reportLines, 0, diagReport.bitCounts, nil, nil}
	diagReport.rootNode.buildAllBelow()
	// After the above we now have a binary tree accessible via the rootNode.
}

func (diagReport *DiagReport) calculateGammaRate() {
	gammaRate := 0b0
	numBits := len(diagReport.bitCounts)
	reportLength := len(diagReport.reportLines)

	// Note that bitCounts[0] holds the MSB!
	for i := 0; i < numBits; i++ {
		if diagReport.bitCounts[i] > reportLength/2 {
			gammaRate = gammaRate | 0b1<<(numBits-i-1)
		}
	}
	diagReport.gammaRate = gammaRate
}

func (diagReport *DiagReport) calculateEpsilonRate() {
	numBits := len(diagReport.bitCounts)
	epsilonRate := diagReport.gammaRate ^ utils.GetMask(numBits)
	diagReport.epsilonRate = epsilonRate
}

func (diagReport DiagReport) getPowerConsumption() int {
	return diagReport.gammaRate * diagReport.epsilonRate
}

func main() {
	diagReport := DiagReport{}
	diagReport.initFromFile(os.Args[1])

	diagReport.calculateGammaRate()
	fmt.Printf("Gamma rate: %d\n", diagReport.epsilonRate)

	diagReport.calculateEpsilonRate()
	fmt.Printf("Epsilon rate: %d\n", diagReport.epsilonRate)

	fmt.Printf("Power consumption: %d\n", diagReport.getPowerConsumption())

	oxygenValue := diagReport.rootNode.search(OXYGEN)[0]
	co2Value := diagReport.rootNode.search(CO2)[0]
	fmt.Println(oxygenValue)
	fmt.Println(co2Value)

	oxygenRating, err := strconv.ParseInt(oxygenValue, 2, 64)
	utils.Check(err)

	co2Rating, err := strconv.ParseInt(co2Value, 2, 64)
	utils.Check(err)

	lifeSupportRating := oxygenRating * co2Rating

	fmt.Printf("Oxygen Rating: %d\nCO2 Rating: %d\nLife Support: %d\n", oxygenRating, co2Rating, lifeSupportRating)
}
