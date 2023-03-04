package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")
	for _, line := range lines {
		fmt.Println(GetMarkerEnd(line))
	}
}

func GetMarkerEnd(line string) int {
	series := []rune(line)
	markerEnd := 0

	for i := 0; ; i++ {
		pool := series[i : i+4]
		duplicatePoint := GetDuplicatePoint(pool)

		if duplicatePoint == 0 {
			markerEnd = i + 3
			break
		}

		// We want to preceed from the position behind the first element of the duplicate.
		i += duplicatePoint - 1
		if i+4 > len(series) {
			return 0
		}
	}

	return markerEnd + 1
}

func GetDuplicatePoint(pool []rune) int {
	visited := make(map[rune]int, 0)
	for i := 0; i < len(pool); i++ {
		if visited[pool[i]] > 0 {
			return visited[pool[i]]
		} else {
			visited[pool[i]] = i + 1
		}
	}
	return 0
}

func ReadLines(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(content), "\n")
}
