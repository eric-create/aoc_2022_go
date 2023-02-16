package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func main() {
	lines, err := read_lines("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	pairs := getPairs(lines)

	sum := 0
	for _, pair := range pairs {
		// fmt.Println(pair[0].String(), pair[1].String())
		if pair[0].IsSubset(pair[1]) || pair[1].IsSubset(pair[0]) {
			sum++
		}
	}

	fmt.Println(sum)
}

func getPairs(lines []string) [][2]mapset.Set[int] {

	setPairs := [][2]mapset.Set[int]{}

	for _, line := range lines {
		var pair []string = strings.Split(line, ",")
		setPairs = append(setPairs, [2]mapset.Set[int]{
			getSectionSet(pair[0]),
			getSectionSet(pair[1])})
	}

	return setPairs
}

func getSectionSet(sections string) mapset.Set[int] {
	var bounds []string = strings.Split(sections, "-")
	start, err := strconv.Atoi(bounds[0])
	if err != nil {
		log.Fatal("failed to section start.")
	}
	end, err := strconv.Atoi(bounds[1])
	if err != nil {
		log.Fatal("failed to section end.")
	}
	set := mapset.NewSet[int]()

	for i := start; i <= end; i++ {
		set.Add(i)
	}

	return set
}

// func get_contained_pairs() []string {

// }

func read_lines(path string) ([]string, error) {
	if content, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(string(content), "\n"), nil
	}
}
