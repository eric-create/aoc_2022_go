package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Knot struct {
	Label  string
	Height rune
	// Up, Right, Left, Down
	Neighbors  []*Knot
	Precursor  *Knot
	Distance   int
	Discovered bool
	X          int
	Y          int
}

func main() {
	lines := ReadFile("./input.txt")
	net := Net(lines)
	startKnot := GetStart(net)
	end := Discover(startKnot)

	// PrintNetDistances(net)

	fmt.Println((*end).Distance, (*end).X, (*end).Y)

	chain := GetChain(end)
	os.WriteFile("output.txt", []byte(SelectionString(net, chain)), 0644)
}

func GetChain(end *Knot) []*Knot {
	chain := []*Knot{}
	cursor := end

	for cursor.Precursor != nil {
		chain = append(chain, cursor.Precursor)
		cursor = cursor.Precursor
	}

	return chain
}

func Discover(knot *Knot) *Knot {
	queue := []*Knot{}
	knot.Discovered = true
	knot.Distance = 0
	return _Discover(knot, queue)
}

func _Discover(knot *Knot, queue []*Knot) *Knot {
	fmt.Println("Discovering", knot.X, knot.Y)
	if knot.Label == "end" {
		return knot
	}

	for _, neighbor := range knot.Neighbors {
		if neighbor != nil {
			if !neighbor.Discovered && neighbor.Height <= knot.Height+1 {
				neighbor.Discovered = true
				neighbor.Distance = knot.Distance + 1
				neighbor.Precursor = knot

				enqueue(&queue, neighbor)
			}
		}
	}

	if len(queue) > 0 {
		next := dequeue(&queue)
		return _Discover(next, queue)
	} else {
		return nil
	}
}

func dequeue(queue *[]*Knot) *Knot {
	knot := (*queue)[0]
	*queue = (*queue)[1:]
	return knot
}

func enqueue(queue *[]*Knot, knot *Knot) {
	*queue = append(*queue, knot)
}

func SelectionString(net [][]*Knot, selection []*Knot) string {
	selectionString := ""

	for _, row := range net {
		for _, knot := range row {
			if Contains(&selection, knot) {
				selectionString += " "
			} else {
				selectionString += string(knot.Height)
			}
		}
		selectionString += "\n"
	}
	return selectionString
}

func GetStart(net [][]*Knot) *Knot {
	for y := 0; y < len(net)-1; y++ {
		for x := 0; x < len((net)[0])-1; x++ {
			if net[y][x].Label == "start" {
				return net[y][x]
			}
		}
	}
	return nil
}

func GetEnd(net [][]*Knot) (x, y int) {
	for y := 0; y < len(net)-1; y++ {
		for x := 0; x < len((net)[0])-1; x++ {
			if (net)[y][x].Label == "end" {
				return x, y
			}
		}
	}
	return -1, -1
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
			fmt.Print(string(knot.Height))
		}
		fmt.Println()
	}
}

func PrintNetDistances(net [][]*Knot) {
	for _, knots := range net {
		for _, knot := range knots {
			if len(strconv.Itoa(knot.Distance)) < 2 {
				fmt.Print("  ", knot.Distance, ";")
			} else {
				fmt.Print(" ", knot.Distance, ";")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func Net(lines []string) [][]*Knot {
	net := [][]*Knot{}

	for y, line := range lines {
		net = append(net, []*Knot{})

		for x, char := range line {
			letter := string(char)
			var knot *Knot = nil

			if letter == "E" {
				knot = &Knot{Label: "end", Height: 'z'}
			} else if letter == "S" {
				knot = &Knot{Label: "start", Height: 'a'}
			} else {
				knot = &Knot{Height: rune(letter[0])}
			}

			knot.Neighbors = []*Knot{nil, nil, nil, nil}
			knot.Y = y
			knot.X = x
			net[y] = append(net[y], knot)
		}
	}

	ProcessNet(&net)
	return net
}

func ProcessNet(net *[][]*Knot) {
	for y, row := range *net {
		for x, knot := range row {

			// Up
			if y > 0 {
				(*knot).Neighbors[0] = (*net)[y-1][x]
			}

			// Right
			if x < len((*net)[0])-1 {
				(*knot).Neighbors[1] = (*net)[y][x+1]
			}

			// Down
			if y < len(*net)-1 {
				(*knot).Neighbors[2] = (*net)[y+1][x]
			}

			// Left
			if x > 0 {
				(*knot).Neighbors[3] = (*net)[y][x-1]
			}
		}
	}
}

func ReadFile(path string) []string {
	text, _ := os.ReadFile(path)
	return strings.Split(string(text), "\n")
}
