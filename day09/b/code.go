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

	head := Knot{11, 5, nil, nil, [][2]int{}, "H"} // 11, 5
	head.AddKnots(8, 11, 5)

	for _, movement := range movements {
		Move(&head, movement)
	}

	tail := head.AbsoluteTail()
	fmt.Println(tail.Name, len(tail.History))
}

func Move(knot *Knot, movement Movement) {
	// Print(*knot)
	knot.ExtendHistory()

	if movement.Distance > 0 {
		knot.XPos += movement.XDir
		knot.YPos += movement.YDir

		if knot.Tail != nil {
			xDistance := math.Abs(float64(knot.XPos) - float64(knot.Tail.XPos))
			yDistance := math.Abs(float64(knot.YPos) - float64(knot.Tail.YPos))

			if xDistance >= 2 || yDistance >= 2 {
				Move(knot.Tail, Movement{
					XDir:     Normalize(knot.XPos - knot.Tail.XPos),
					YDir:     Normalize(knot.YPos - knot.Tail.YPos),
					Distance: 1})
			}
		}

		movement.Distance -= 1
		Move(knot, movement)
	}
}

func Print(knot Knot) {
	for y := 20; y >= 0; y-- { // 20 4
		for x := 0; x <= 25; x++ { // 25 5
			if x == 11 && y == 5 {
				fmt.Print("s")
			} else {
				fmt.Print(GetSymbol(x, y, *knot.AbsoluteHead()))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func GetSymbol(xPos int, yPos int, knot Knot) string {
	if knot.XPos == xPos && knot.YPos == yPos {
		return knot.Name
	} else if knot.Tail != nil {
		return GetSymbol(xPos, yPos, *knot.Tail)
	} else {
		return "."
	}
}

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
	Name    string
}

func (knot Knot) AbsoluteHead() *Knot {
	if knot.Head != nil {
		return knot.Head.AbsoluteHead()
	}
	return &knot
}

func (knot Knot) AbsoluteTail() *Knot {
	if knot.Tail != nil {
		return knot.Tail.AbsoluteTail()
	}
	return &knot
}

func (knot *Knot) AddKnots(knotCount int, xPos int, yPos int) {
	if knotCount >= 0 {
		knot.Tail = &Knot{xPos, yPos, knot, nil, [][2]int{}, strconv.Itoa(9 - knotCount)}
		knot.Tail.AddKnots(knotCount-1, xPos, yPos)
	}
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
