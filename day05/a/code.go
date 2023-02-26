package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := getLines("./input.txt")
	stacks := getstacks(lines)
	operations := getOperations(lines)

	for _, operation := range operations {
		operate(&stacks, operation)
	}

	printTopping(stacks)
}

func printTopping(stacks [][]rune) {
	for _, stack := range stacks {
		fmt.Print(string(stack[len(stack)-1]))
	}
	fmt.Println()
}

func operate(stacks *[][]rune, operation Operation) {
	for i := operation.Number; i > 0; i-- {

		sourceStack := (*stacks)[operation.Source]
		targetStack := (*stacks)[operation.Target]
		crate := sourceStack[len(sourceStack)-1]

		(*stacks)[operation.Target] = append(targetStack, crate)
		(*stacks)[operation.Source] = sourceStack[:len(sourceStack)-1]
	}
}

type Operation struct {
	Number int
	Source int
	Target int
}

func NewOperation(number int, source int, target int) Operation {
	return Operation{
		Number: number,
		Source: source - 1,
		Target: target - 1}
}

func getOperation(line string) Operation {
	regex := regexp.MustCompile(`move (?P<number>\d+) from (?P<source>\d+) to (?P<target>\d+)`)
	matches := regex.FindStringSubmatch(line)

	number, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal("failed to get the number")
	}

	source, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Fatal("failed to get the source")
	}

	target, err := strconv.Atoi(matches[3])
	if err != nil {
		log.Fatal("failed to get the target")
	}

	return NewOperation(number, source, target)
}

func getOperations(lines []string) []Operation {
	operationLines := getOperationLines(lines)
	operations := []Operation{}

	for _, line := range operationLines {
		operations = append(operations, getOperation(line))
	}

	return operations
}

func getOperationLines(lines []string) []string {
	operationLines := []string{}

	for _, line := range lines {
		if strings.HasPrefix(line, "move") {
			operationLines = append(operationLines, line)
		}
	}

	return operationLines
}

func getstacks(lines []string) [][]rune {
	crateLines := getCrateLines(lines)
	stacks := make([][]rune, numOfStacks(crateLines))

	for i := len(crateLines) - 1; i >= 0; i-- {
		for j, char := range crateLines[i] {
			if char != ' ' {
				stacks[j] = append(stacks[j], char)
			}
		}
	}

	return stacks
}

func numOfStacks(crateLines []string) int {
	stacksSum := 0

	for _, line := range crateLines {
		if len(line) > stacksSum {
			stacksSum = len(line)
		}
	}

	return stacksSum
}

func getCrateLines(lines []string) []string {
	crateLines := make([]string, len(lines))

	for i, line := range lines {
		if strings.Contains(line, "[") {
			for j := 1; j < len(line); j += 4 {
				crateLines[i] += string(line[j])
			}
		}
	}

	return crateLines
}

func getLines(path string) []string {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(file), "\n")
}
