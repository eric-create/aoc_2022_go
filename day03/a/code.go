package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	path := "./input.txt"
	lines, err := read(path)

	if err != nil {
		log.Fatal("Failed to read")
	}

	duplicates := []string{}
	priorities := []int{}

	for _, line := range lines {
		first, second := get_compartements(line)

		duplicate, err := get_duplicate(first, second)
		if err != nil {
			log.Fatal(err)
		}

		duplicates = append(duplicates, duplicate)
		priorities = append(priorities, get_priority(duplicate))
	}

	sum := 0

	for _, priority := range priorities {
		sum += priority
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

	fmt.Println(priority)
	return priority
}

func get_duplicate(first string, second string) (string, error) {
	for _, char := range first {
		if strings.Contains(second, string(char)) {
			return string(char), nil
		}
	}
	return "", errors.New(fmt.Sprintf("Could not find a duplicate in %v %v", first, second))
}

func get_compartements(line string) (string, string) {
	return line[:len(line)/2], line[len(line)/2:]
}

func read(path string) ([]string, error) {
	if content, err := ioutil.ReadFile(path); err == nil {
		return strings.Split(string(content), "\n"), nil
	} else {
		return nil, err
	}
}
