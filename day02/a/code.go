package main

import (
	"fmt"
	"os"
	"strings"
)

type Match struct {
	Call, Response string
}

func main() {
	lines, err := getLines("./input.txt")
	exitOnError(err)
	matches := getMatches(lines)
	normalize(&matches)
	scores := getScores(matches)

	for _, score := range scores {
		fmt.Println(score)
	}
}

func getScores(matches []Match) []int {
	scores := []int{}

	for _, match := range matches {
		score := 0
		score += getShapeScore(match)
		score += getMatchScore(match)
		scores = append(scores, score)
	}

	return scores
}

type Symbol struct {
	Above   *Symbol
	Beneath *Symbol
}

// Returns three instances of Symbol: "rock", "paper", and "scissors".
func getSymbols() (rock, paper, scissors Symbol) {
	rock = Symbol{}
	paper = Symbol{}
	scissors = Symbol{}

	rock.Above = &scissors
	rock.Beneath = &paper
	paper.Above = &rock
	paper.Beneath = &scissors
	scissors.Above = &paper
	scissors.Beneath = &rock

	return rock, paper, scissors
}

// Returns the represented symbol for a given letter.
func getSymbol(letter string) Symbol {
	rock, paper, scissors := getSymbols()
	var symbol Symbol

	switch letter {
	case "A":
		symbol = rock
	case "B":
		symbol = paper
	case "C":
		symbol = scissors
	}

	return symbol
}

func getMatchScore(match Match) int {
	call := getSymbol(match.Call)
	response := getSymbol(match.Response)
	matchScore := 0

	switch response {
	case *call.Beneath:
		matchScore = 6
	case call:
		matchScore = 3
	case *call.Above:
		matchScore = 0
	}
	return matchScore
}

func getShapeScore(match Match) int {
	switch response := match.Response; response {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	}
	return 0
}

func normalize(matches *[]Match) {
	for i := range *matches {
		switch response := &((*matches)[i].Response); *response {
		case "X":
			*response = "A"
		case "Y":
			*response = "B"
		case "Z":
			*response = "C"
		}
	}
}

func getMatches(lines []string) []Match {
	matches := []Match{}
	for _, line := range lines {
		lineTuple := strings.Split(line, " ")
		matches = append(matches, Match{Call: lineTuple[0], Response: lineTuple[1]})
	}
	return matches
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getLines(path string) ([]string, error) {
	if input, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(string(input), "\n"), nil
	}
}
