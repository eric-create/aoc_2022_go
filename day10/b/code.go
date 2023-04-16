package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")
	operations := Operations(lines)
	cpu := CPU{X: 1, History: []int{}}

	for _, operation := range operations {
		cpu.Register(operation)
	}

	PrintScreen(cpu.History)
}

func PrintScreen(history []int) {
	for i, spritePosition := range history {
		pixelPosition := i % 40
		if IsSpriteVisible(pixelPosition, spritePosition) {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
		if pixelPosition == 39 {
			fmt.Println()
		}
	}
}

func IsSpriteVisible(pixelPosition int, spritePosition int) bool {
	if pixelPosition == spritePosition ||
		pixelPosition == spritePosition+1 ||
		pixelPosition == spritePosition-1 {
		return true
	}
	return false
}

type CPU struct {
	X       int
	History []int
}

func (cpu *CPU) Register(operation Operation) {
	for i := 0; i < operation.Duration; i++ {
		(*cpu).History = append((*cpu).History, (*cpu).X)
	}
	(*cpu).X += operation.Impact
}

func (cpu *CPU) GetCycle(cycle int) int {
	return (*cpu).History[cycle-1]
}

func (cpu *CPU) GetSignalStrength(cycle int) int {
	return cycle * (*cpu).GetCycle(cycle)
}

func (cpu *CPU) GetSignalsSum(cycles []int) int {
	sum := 0
	for _, cycle := range cycles {
		sum += cpu.GetSignalStrength(cycle)
	}
	return sum
}

func Operations(lines []string) []Operation {
	operations := []Operation{}

	for _, line := range lines {
		operations = append(operations, GetOperation(line))
	}
	return operations
}

func GetOperation(line string) Operation {
	parts := strings.Split(line, " ")

	if len(parts) > 1 {
		impact, _ := strconv.Atoi(parts[1])
		return Operation{Duration: 2, Impact: impact}

	} else {
		return Operation{Duration: 1, Impact: 0}
	}
}

type Operation struct {
	Duration int
	Impact   int
}

func ReadLines(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
