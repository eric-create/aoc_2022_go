package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")
	// Appending an empty string to the array is important, so that the last monkey defined
	// in the text file is being processed.
	lines = append(lines, "")
	monkeys := GetMonkeys(lines)
	rounds := 20

	for i := 1; i <= rounds; i++ {
		for _, monkey := range monkeys {
			for len(monkey.Items) > 0 {
				item := monkey.PopItem()
				item = monkey.Inspect(item)
				test := monkey.Test(item)
				peer := monkey.GetPeer(test, monkeys)
				monkey.Throw(item, peer)
			}
		}
	}

	for i, monkey := range monkeys {
		fmt.Println("Monkey", i, "Items", monkey.Items)
	}
	fmt.Println()
	business := GetMonkeyBusiness(monkeys)
	fmt.Println()
	fmt.Println("Monkey Business", business)
}

func GetMonkeyBusiness(monkeys []*Monkey) int {
	businesses := GetBusinesses(monkeys)
	monkeyCount := len(monkeys)
	return businesses[monkeyCount-1] * businesses[monkeyCount-2]
}

func GetBusinesses(monkeys []*Monkey) []int {
	businesses := []int{}

	for i, monkey := range monkeys {
		businesses = append(businesses, monkey.Business)

		fmt.Println("Business of monkey", i, "is: ", businesses[i])
	}

	sort.Ints(businesses)

	return businesses
}

func GetMonkeys(lines []string) []*Monkey {
	monkeys := []*Monkey{}
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

func GetMonkey(lines []string) *Monkey {

	monkey := &Monkey{Business: 0}

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
	operator := operationSequence[0]
	operand := operationSequence[2:]

	// Looking at the input text then it turns out that the only mathematical operations
	// applied are multiplication and addition.
	if operator == '+' {
		return func(old int) int {
			if operand == "old" {
				return old + old
			} else {
				constant, _ := strconv.Atoi(operand)
				return old + constant
			}
		}
	} else {
		return func(old int) int {
			if operand == "old" {
				return old * old
			} else {
				constant, _ := strconv.Atoi(operand)
				return old * constant
			}
		}
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
	Business       int
}

func (monkey *Monkey) PopItem() int {
	monkey.Business++
	item := monkey.Items[0]
	monkey.Items = monkey.Items[1:]
	return item
}

func (monkey *Monkey) Inspect(item int) int {
	item = monkey.operation(item)
	item = item / 3
	return item
}

func (monkey *Monkey) Test(item int) bool {
	return item%monkey.TestingDivisor == 0
}

func (monkey *Monkey) GetPeer(test bool, monkeys []*Monkey) *Monkey {
	if test {
		return monkeys[monkey.PeerTrue]
	} else {
		return monkeys[monkey.PeerFalse]
	}
}

func (monkey *Monkey) Throw(item int, peer *Monkey) {
	peer.Items = append(peer.Items, item)
}

func ReadLines(path string) []string {
	bytes, _ := os.ReadFile(path)
	content := string(bytes)
	lines := strings.Split(content, "\n")
	return lines
}
