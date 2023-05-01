package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")
	// Appending an empty string to the array is important, so that the last monkey defined
	// in the text file is being processed.
	lines = append(lines, "")
	GetMonkeys(lines)
}

func GetMonkeys(lines []string) []Monkey {
	monkeys := []Monkey{}
	monkeyLines := []string{}

	for _, line := range lines {
		if line == "" {
			monkeys = append(monkeys, GetMonkey(monkeyLines))
			monkeyLines = []string{}
		} else {
			monkeyLines = append(monkeyLines, line)
		}
	}
	return monkeys
}

func GetMonkey(lines []string) Monkey {

	monkey := Monkey{}

	for _, line := range lines {
		if strings.Contains(line, "Starting items") {
			monkey.Items = GetStartingItems(line)
		} else if strings.Contains(line, "Operation") {
			monkey.operation = GetOperation(line)
		} else if strings.Contains(line, "Test") {
			monkey.TestingDivisor = GetDivisor(line)
		} else if strings.Contains(line, "true") {
			monkey.PeerTrue = GetPeer(line, "true")
		} else if strings.Contains(line, "false") {
			monkey.PeerFalse = GetPeer(line, "false")
		}
	}

	fmt.Println(monkey.Items, monkey.TestingDivisor, monkey.PeerTrue, monkey.PeerFalse)
	fmt.Println()
	return monkey
}

func GetStartingItems(line string) []int {
	items := []int{}
	itemsSequence := string(line[18:])
	for _, itemString := range strings.Split(itemsSequence, ", ") {
		item, _ := strconv.Atoi(itemString)
		items = append(items, item)
	}
	return items
}

func GetOperation(line string) Operation {
	operationSequence := string(line[23:])
	operand, _ := strconv.Atoi(operationSequence[2:])

	// Looking at the input text then it turns out that the only mathematical operations
	// applied are multiplication and addition.
	if operationSequence[0] == '+' {
		return func(old int) int { return old + operand }
	} else {
		return func(old int) int { return old * operand }
	}
}

func GetDivisor(line string) int {
	divisorSequence := string(line[21:])
	divisor, _ := strconv.Atoi(divisorSequence)
	return divisor
}

func GetPeer(line string, result string) int {
	sequenceStart := 30
	if result == "true" {
		sequenceStart = 29
	}
	peerSequence := line[sequenceStart:]
	peer, _ := strconv.Atoi(peerSequence)
	return peer
}

type Operation func(old int) int

type Monkey struct {
	Items          []int
	operation      Operation
	TestingDivisor int
	PeerTrue       int
	PeerFalse      int
}

func (monkey *Monkey) PopItem() int {
	item := monkey.Items[0]
	monkey.Items = monkey.Items[1:]
	return item
}

func (monkey *Monkey) Inspect() {
	item := monkey.PopItem()
	item = monkey.operation(item)
}

func ReadLines(path string) []string {
	bytes, _ := os.ReadFile(path)
	content := string(bytes)
	lines := strings.Split(content, "\n")
	return lines
}
