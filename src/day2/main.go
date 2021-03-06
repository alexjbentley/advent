package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"utils"
)

const (
	FORWARD = "forward"
	UP      = "up"
	DOWN    = "down"
)

type Craft struct {
	x, y, aim int
}

type Instruction struct {
	mode  string
	value int
}

func (craft Craft) GetEndPositionValue() int {
	return craft.x * craft.y
}

func (craft *Craft) MoveX(distance int) {
	craft.x += distance
}

func (craft *Craft) MoveY(distance int) {
	craft.y += distance
}

func (craft *Craft) MoveForward(distance int) {
	craft.MoveX(distance)
	craft.MoveY(distance * craft.aim)
}

func (craft *Craft) AimAdjust(adjustment int) {
	craft.aim += adjustment
}

func (craft *Craft) FollowInstruction(instruction Instruction) {
	switch instruction.mode {
	case FORWARD:
		craft.MoveForward(instruction.value)

	case UP:
		craft.AimAdjust(-1 * instruction.value)

	case DOWN:
		craft.AimAdjust(instruction.value)

	default:
		panic("Unrecognised instruction!")
	}
}

func (craft *Craft) FollowRoute(routeFile string) {
	data, err := os.ReadFile(routeFile)
	utils.Check(err)

	instructionStrings := strings.Split(string(data), "\n")

	for _, instructionString := range instructionStrings {
		tokens := strings.Split(instructionString, " ")
		mode := tokens[0]
		value, err := strconv.Atoi(tokens[1])
		utils.Check(err)

		instruction := Instruction{mode, value}
		craft.FollowInstruction(instruction)
	}
}

func main() {
	submarine := Craft{0, 0, 0} // Submarine starts at the origin with neutral aim
	submarine.FollowRoute("files/route_plan.txt")
	fmt.Println(submarine.GetEndPositionValue())
}
