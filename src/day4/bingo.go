package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"utils"
)

type BingoEntry struct {
	value   int
	checked bool
}

func (entry *BingoEntry) check(called int) bool {
	if entry.value == called {
		entry.checked = true
		return true
	}
	return false // Cannot simply return entry.checked, as it may have been checked in an earlier round
}

type BingoLine struct {
	entries []BingoEntry
	checked int
}

func (line *BingoLine) call(called int) bool {
	for i := 0; i < len(line.entries); i++ {
		if line.entries[i].check(called) {
			line.checked++
		}
	}
	return line.checked >= len(line.entries)
}

type BingoCard struct {
	lines []BingoLine
	score int
}

// From the horizontal lines we read in, build the corresponding vertical lines.
func (card *BingoCard) buildVerticalLines(gridSize int) {
	newLines := make([]BingoLine, gridSize)
	for i := 0; i < len(newLines); i++ {
		newLines[i].entries = make([]BingoEntry, gridSize)
	}
	for i, line := range card.lines {
		for j, entry := range line.entries {
			newLines[j].entries[i] = entry
		}
	}
	card.lines = append(card.lines, newLines...)
}

func (card *BingoCard) call(called int) bool {
	for i := 0; i < len(card.lines); i++ {
		if card.lines[i].call(called) {
			// Bingo
			fmt.Printf("Calling %d took card to %d\n", called, card.lines[i].checked)
			return true
		} else {
			fmt.Printf("Calling %d took card to %d\n", called, card.lines[i].checked)
		}
	}
	return false
}

func (card *BingoCard) calculateScore() (score int) {
	// Half of the lines contain all the values for scoring purposes
	score = 0
	for i := float64(0); i < float64(len(card.lines))/2; i++ {
		// For half of the lines (we doubled up)
		for _, entry := range card.lines[int(i)].entries {
			// Sum the unmarked entries
			if !entry.checked {
				fmt.Printf("Adding %d\n", entry.value)
				score += entry.value
			} else {
				fmt.Printf("Not adding %d, entry is marked.\n", entry.value)
			}
		}
	}
	card.score = score
	return
}

type BingoGame struct {
	cards       []BingoCard
	calls       []int
	winningCard BingoCard
	lastCall    int
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

func (game *BingoGame) loadBingoCards(input []string) {
	gridSize := len(strings.Split(input[0], " "))
	cards := make([]BingoCard, len(input)/gridSize)
	gridCounter := 0 // All cards are together need to track where we are.
	cardCounter := 0 // Not for cheating purposes.
	bingoLines := make([]BingoLine, gridSize)
	for _, line := range input {
		// For each line

		bingoLines[gridCounter].entries = make([]BingoEntry, gridSize)
		tokens := strings.Split(line, " ")
		for j, token := range tokens {
			// For each value
			parsedToken, err := strconv.ParseInt(token, 10, 32)
			utils.Check(err)
			bingoLines[gridCounter].entries[j] = BingoEntry{int(parsedToken), false}
		}
		gridCounter++
		if gridCounter == gridSize {
			// We have read in a full card, do some cleanup and read the vertical lines.
			gridCounter = 0
			cards[cardCounter] = BingoCard{bingoLines, 0}
			cards[cardCounter].buildVerticalLines(gridSize)
			cardCounter++
			bingoLines = make([]BingoLine, gridSize) // Reset the working set of lines.
		}
	}
	game.cards = cards
}

func (game *BingoGame) call(called int) bool {
	for i := 0; i < len(game.cards); i++ {
		if game.cards[i].call(called) {
			// This is the winning card
			fmt.Println("Winning card found!")
			game.winningCard = game.cards[i]
			return true
		}
	}
	return false
}

func (game *BingoGame) play() {
	for _, call := range game.calls {
		fmt.Printf("Calling: %d\n", call)
		if game.call(call) {
			// Bingo has been found, winningCard will be set
			// Calculate the score
			game.winningCard.calculateScore()
			// Store the final call
			game.lastCall = call
			return
		}
	}
}

func preProcessData(filename string) []string {
	data, err := os.ReadFile(filename)
	utils.Check(err)

	processed := strings.ReplaceAll(string(data), "\n\n", "\n")
	processed = strings.ReplaceAll(processed, "  ", " ")

	processedLines := strings.Split(processed, "\n")
	for i := 0; i < len(processedLines); i++ {
		processedLines[i] = strings.Trim(processedLines[i], " ")
		processedLines[i] = strings.TrimLeft(processedLines[i], " ")
	}
	return processedLines
}

func main() {
	var bingoGame BingoGame
	processedData := preProcessData(os.Args[1])
	callString := processedData[0]
	cardLines := processedData[1:]
	bingoGame.loadBingoCalls(callString)
	bingoGame.loadBingoCards(cardLines)
	bingoGame.play()
	fmt.Printf("Winning card: %+v\n\n\n", bingoGame.winningCard)
	fmt.Printf("Last Call: %d\n", bingoGame.lastCall)
	fmt.Printf("Card score: %d\n", bingoGame.winningCard.score)
	fmt.Printf("Final score: %d\n", bingoGame.lastCall*bingoGame.winningCard.score)
}
