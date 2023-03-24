package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")
	movements := Movements(lines)

	tail := Knot{0, 0, nil, nil, [][2]int{}}
	head := Knot{0, 0, nil, &tail, [][2]int{}}
	tail.Head = &head

	for _, movement := range movements {
		Move(&head, movement)
	}

	fmt.Println(len(tail.History))
}

func Move(knot *Knot, movement Movement) {
	knot.ExtendHistory()

	if movement.Distance > 0 {
		knot.XPos += movement.XDir
		knot.YPos += movement.YDir

		if knot.Tail != nil {
			xDistance := math.Abs(float64(knot.XPos) - float64(knot.Tail.XPos))
			yDistance := math.Abs(float64(knot.YPos) - float64(knot.Tail.YPos))

			// Tail follows in horizontal or vertical direction.
			if xDistance+yDistance == 2 && (xDistance == 2 || yDistance == 2) {
				Move(knot.Tail, Movement{movement.XDir, movement.YDir, 1})
			}

			// Tail follows in diagonal direction.
			if xDistance+yDistance == 3 {
				Move(knot.Tail, Movement{
					XDir:     Normalize(knot.XPos - knot.Tail.XPos),
					YDir:     Normalize(knot.YPos - knot.Tail.YPos),
					Distance: 1})
			}
		}

		// Print(*head)
		movement.Distance -= 1
		Move(knot, movement)
	}
}

// // Only works for the example
// func Print(knot Knot) {
// 	head := knot.GetCapital()

// 	for y := 4; y >= 0; y-- {
// 		for x := 0; x <= 5; x++ {
// 			if head.XPos == x && head.YPos == y {
// 				fmt.Print("H")
// 			} else if head.Tail.XPos == x && head.Tail.YPos == y {
// 				fmt.Print("T")
// 			} else {
// 				fmt.Print(".")
// 			}
// 		}
// 		fmt.Println()
// 	}
// 	fmt.Println()
// }

func Normalize(i int) int {
	if i < 0 {
		return -1
	} else if i > 1 {
		return 1
	} else {
		// i is 0.
		return i
	}
}

type Knot struct {
	XPos    int
	YPos    int
	Head    *Knot
	Tail    *Knot
	History [][2]int
}

func (knot Knot) GetCapital() *Knot {
	if knot.Head != nil {
		return knot.Head.GetCapital()
	}
	return &knot
}

func (knot *Knot) ExtendHistory() {
	for _, position := range knot.History {
		if position[0] == knot.XPos && position[1] == knot.YPos {
			return
		}
	}
	knot.History = append(knot.History, [2]int{knot.XPos, knot.YPos})
}

func Movements(lines []string) []Movement {
	movements := []Movement{}

	for _, line := range lines {
		instructions := strings.Split(line, " ")
		xDir, yDir := TranslateDirection(instructions[0])
		distance, _ := strconv.Atoi(instructions[1])
		movements = append(movements, Movement{xDir, yDir, distance})
	}

	return movements
}

func TranslateDirection(direction string) (int, int) {
	switch direction {
	case "U":
		return 0, 1
	case "D":
		return 0, -1
	case "L":
		return -1, 0
	case "R":
		return 1, 0
	}
	log.Fatal(direction)
	return 0, 0
}

func ReadLines(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	return lines
}

type Movement struct {
	XDir     int
	YDir     int
	Distance int
}
