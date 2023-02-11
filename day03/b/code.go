package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func main() {
	path := "./input.txt"
	lines, err := read(path)

	if err != nil {
		log.Fatal("Failed to read")
	}

	badges := []rune{}

	for i := 0; i < len(lines); i += 3 {
		s := mapset.NewSet([]rune(lines[i])...).
			Intersect(mapset.NewSet([]rune(lines[i+1])...)).
			Intersect(mapset.NewSet([]rune(lines[i+2])...))

		badges = append(badges, s.ToSlice()[0])
	}

	sum := 0
	for _, badge := range badges {
		sum += get_priority(string(badge))
	}
	fmt.Println("Sum:", sum)
}

func get_priority(duplicate string) int {
	priority := 0

	if int(duplicate[0]) < 97 {
		priority = int(duplicate[0]) - 38
	} else {
		priority = int(duplicate[0]) - 96
	}

	return priority
}

func read(path string) ([]string, error) {
	if content, err := os.ReadFile(path); err == nil {
		return strings.Split(string(content), "\n"), nil
	} else {
		return nil, err
	}
}
