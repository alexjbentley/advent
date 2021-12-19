package main

import (
	"fmt"
	"os"
	"strings"
	"utils"
)

var ZERO = string("0")
var ONE = string("1")

type Node struct {
	// TODO Need to have our own bitCounts for the values we possess. Do away with numBits, we can get it from our own
	// bit counts once we have them.
	values []string

	// The particular "bit" this node is concerned with. Analagous with depth.
	index, numBits int
	left           *Node
	right          *Node
}

// Recursively build all nodes below this one, based on the values
func (node *Node) buildAllBelow() {
	// TODO we shouldn't be storing any values in the non-leaf nodes. Memory efficiency!
	if node.index == node.numBits || len(node.values) == 1 {
		// We're the leaf of the tree, stop.
		return
	} else {
		// Make left Node (0):
		// Decide which values should go in it first
		// TODO: make these two into a func
		var workingValues []string
		var bit string
		for _, value := range node.values {
			bit = strings.Split(value, "")[node.index]
			if bit == ZERO {
				workingValues = append(workingValues, value)
			}
		}

		fmt.Printf("Creating left node with values:\n")
		fmt.Println(workingValues)
		if len(workingValues) > 0 {
			node.left = &Node{workingValues, node.index + 1, node.numBits, nil, nil}
			if len(workingValues) > 1 {
				node.left.buildAllBelow()
			}
		}

		workingValues = []string{}

		// Make right Node (1):
		// Decide which values should go in it first
		for _, value := range node.values {
			bit = strings.Split(value, "")[node.index]
			if bit == ONE {
				workingValues = append(workingValues, value)
			}
		}

		fmt.Println("Creating right node with values:")
		fmt.Println(workingValues)
		if len(workingValues) > 0 {
			node.right = &Node{workingValues, node.index + 1, node.numBits, nil, nil}
			if len(workingValues) > 1 {
				node.right.buildAllBelow()
			}
		}
	}
}

func (node Node) searchWithSequence(sequence []string) []string {
	// Escape clause first!
	fmt.Println("Node with following being searched: ")
	fmt.Println(node.values)
	// if len(node.values) == 1 || node.index == node.numBits-1 {
	if len(node.values) == 1 {
		return node.values
	}
	if sequence[node.index] == ONE {
		return node.right.searchWithSequence(sequence)
	} else {
		return node.left.searchWithSequence(sequence)
	}
}

type DiagReport struct {
	gammaRate, epsilonRate, oxygenRating, co2Rating, lifeSupportRating int
	reportLines                                                        []string
	rootNode                                                           *Node
	bitCounts                                                          []int
}

func (diagReport *DiagReport) initFromFile(filename string) {
	data, err := os.ReadFile(filename)
	utils.Check(err)
	diagReport.reportLines = strings.Split(string(data), "\n")
	numBits := len(diagReport.reportLines[0])

	// TODO Move this after the bitCounts calculations, then pass bitCounts in. Will need recalculating at each node.
	diagReport.rootNode = &Node{diagReport.reportLines, 0, numBits, nil, nil}
	diagReport.rootNode.buildAllBelow()
	// After the above we now have a binary tree accessible via the rootNode.

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

func (diagReport DiagReport) constructSearchSequences() (oxygenSequence []string, co2Sequence []string) {
	oxygenSequence, co2Sequence = make([]string, 0), make([]string, 0)
	reportLength := len(diagReport.reportLines)
	numBits := len(diagReport.bitCounts)
	for i := 0; i < numBits; i++ {
		if diagReport.bitCounts[i] >= reportLength/2 {
			oxygenSequence = append(oxygenSequence, ONE)
			co2Sequence = append(co2Sequence, ZERO)
		} else {
			oxygenSequence = append(oxygenSequence, ZERO)
			co2Sequence = append(co2Sequence, ONE)
		}
	}
	return
}

func main() {
	diagReport := DiagReport{}
	diagReport.initFromFile(os.Args[1])

	diagReport.calculateGammaRate()
	fmt.Printf("Gamma rate: %d\n", diagReport.epsilonRate)

	diagReport.calculateEpsilonRate()
	fmt.Printf("Epsilon rate: %d\n", diagReport.epsilonRate)

	fmt.Printf("Power consumption: %d\n", diagReport.getPowerConsumption())

	oxygenSearchSequence, co2SearchSequence := diagReport.constructSearchSequences()
	fmt.Printf("About to search with sequences:\nOxygen: %v\nCO2: %v\n", oxygenSearchSequence, co2SearchSequence)

	// Hopefully both of these have returned only one value. Beware of duplicates?
	oxygenValues := diagReport.rootNode.searchWithSequence(oxygenSearchSequence)
	co2Values := diagReport.rootNode.searchWithSequence(co2SearchSequence)
	fmt.Println(oxygenValues)
	fmt.Println(co2Values)
}
