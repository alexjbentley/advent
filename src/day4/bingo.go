package main

import (
	"os"
	"strconv"
	"strings"
	"utils"
)

type BingoEntry struct {
	value   int
	checked bool
}

type BingoLine struct {
	values []BingoEntry
}

type BingoCard struct {
	lines []BingoLine
}

type BingoGame struct {
	cards []BingoCard
	calls []int
}

func (game *BingoGame) loadBingoCalls(input string) {
	stringCalls := strings.Split(input, ",")
	calls := make([]int, len(stringCalls))
	for i, call := range stringCalls {
		parsedInt, err := strconv.ParseInt(call, 10, 64)
		utils.Check(err)

		calls[i] = int(parsedInt)

	}
	game.calls = calls
}

func (game *BingoGame) loadBingoCards(input []string) []BingoCard {
	gridSize := len(strings.Split(input[0], " "))
	cards := make([]BingoCard, len(input)/gridSize)
	for i, line := range input {
		// For each line
		// bingoLine :=
	}
}

func preProcessData(filename string) []string {
	data, err := os.ReadFile(filename)
	utils.Check(err)

	processed := strings.ReplaceAll(string(data), "\n\n", "\n")
	processed = strings.ReplaceAll(processed, "  ", " ")

	return strings.Split(processed, "\n")
}

func main() {
	var bingoGame BingoGame
	processedData := preProcessData(os.Args[1])
	callString := processedData[0]
	cardLines := processedData[1:]
	bingoGame.loadBingoCalls(callString)
	bingoGame.loadBingoCards(cardLines)

}
