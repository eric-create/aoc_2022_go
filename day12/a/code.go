package main

import (
	"fmt"
	"os"
	"strings"
)

type Knot struct {
	Label   string
	Height  string
	Neighs  []*Knot
	Ascents []*Knot
}

func main() {
	lines := ReadFile("./input.txt")
	net := Net(lines)
	PrintNet(net)
	x, y := GetEnd(&net)
	zplane := GetPlane(x, y, "z", net)
	PrintNetSelection(net, zplane)
}

func PrintNetSelection(net [][]*Knot, selection []*Knot) {
	for _, row := range net {
		for _, knot := range row {
			if Contains(&selection, knot) {
				fmt.Print(knot.Height)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func GetEnd(net *[][]*Knot) (x, y int) {
	for y := 0; y < len(*net)-1; y++ {
		for x := 0; x < len((*net)[0])-1; x++ {
			if (*net)[y][x].Height == "E" {
				(*net)[y][x].Height = "z"
				return x, y
			}
		}
	}
	return -1, -1
}

func GetPlane(x, y int, height string, net [][]*Knot) []*Knot {
	plane := []*Knot{}
	_GetPlane(x, y, height, net, &plane)
	return plane
}

func _GetPlane(x, y int, height string, net [][]*Knot, plane *[]*Knot) {
	knot := net[y][x]

	if Contains(plane, knot) {
		return
	}
	*plane = append(*plane, knot)

	// Left
	if x > 0 && net[y][x-1].Height == height {
		_GetPlane(x-1, y, height, net, plane)
	}

	// Right
	if x < len(net[0])-1 && net[y][x+1].Height == height {
		_GetPlane(x+1, y, height, net, plane)
	}

	// Up
	if y > 0 && net[y-1][x].Height == height {
		_GetPlane(x, y-1, height, net, plane)
	}

	// Down
	if y < len(net)-1 && net[y+1][x].Height == height {
		_GetPlane(x, y+1, height, net, plane)
	}
}

func Contains[T any](list *[]*T, object *T) bool {
	for _, element := range *list {
		if element == object {
			return true
		}
	}
	return false
}

func PrintNet(net [][]*Knot) {
	for _, knots := range net {
		for _, knot := range knots {
			fmt.Print(knot.Height, "")
		}
		fmt.Println()
	}
}

func Net(lines []string) [][]*Knot {
	knots := [][]*Knot{}
	for y, line := range lines {
		knots = append(knots, []*Knot{})
		for _, char := range line {
			knots[y] = append(knots[y], &Knot{Height: string(char)})
		}
	}
	return knots
}

func ReadFile(path string) []string {
	text, _ := os.ReadFile(path)
	return strings.Split(string(text), "\n")
}
