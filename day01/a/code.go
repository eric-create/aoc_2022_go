package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	lines := getLines("./input.txt")
	bags := getBags(lines)

	sort.Ints(bags)

	// Get the biggest bag.
	fmt.Println(strconv.Itoa(bags[len(bags)-1]))
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getLines(path string) []string {
	input, err := os.ReadFile(path)
	exitOnError(err)

	lines := strings.Split(string(input), "\n")
	return lines
}

func getBags(lines []string) []int {
	bags := []int{}
	var bag int

	for _, line := range lines {
		if line != "" {

			// Get the calories of the food.
			calories, err := strconv.Atoi(line)
			exitOnError(err)

			// Append the calories to the current bag.
			bag += calories

		} else {

			bags = append(bags, bag)
			bag = 0

		}
	}

	bags = append(bags, bag)

	return bags
}
