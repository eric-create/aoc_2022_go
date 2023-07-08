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
	startX, startY := GetStart(net)
	end := Discover(&net, startX, startY)

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

func Discover(net *[][]*Knot, xStart int, yStart int) *Knot {
	knot := (*net)[yStart][xStart]

	var link *Knot = nil
	var end_dummy *Knot = nil
	var end **Knot = &end_dummy

	ascender, descender := DiscoverPlane(knot, 0, end)
	// PrintNetDistances(*net)

	for *end == nil {
		if ascender == nil {
			link = descender
		} else {
			link = ascender
		}
		ascender, descender = DiscoverPlane(link, link.Distance, end)
		// PrintNetDistances(*net)
	}

	return *end
}

func DiscoverPlane(knot *Knot, distance int, end **Knot) (*Knot, *Knot) {
	knot.Distance = distance
	knot.Discovered = true

	var ascender_dummy *Knot = nil
	var descender_dummy *Knot = nil
	var ascender **Knot = &ascender_dummy
	var descender **Knot = &descender_dummy

	queue := []*Knot{}
	_DiscoverPlane(knot, ascender, descender, end, &queue)

	return *ascender, *descender
}

func _DiscoverPlane(knot *Knot, ascender **Knot, descender **Knot, end **Knot, queue *[]*Knot) {

	if *ascender != nil || *end != nil {
		return
	}

	for _, neighbor := range knot.Neighbors {
		if neighbor != nil {

			if !neighbor.Discovered && neighbor.Height == knot.Height {
				neighbor.Precursor = knot
				neighbor.Distance = knot.Distance + 1
				neighbor.Discovered = true

				enqueue(queue, neighbor)

			} else if *ascender == nil && neighbor.Height == knot.Height+1 {
				*ascender = neighbor
				(*ascender).Precursor = knot
				(*ascender).Distance = knot.Distance + 1
				(*ascender).Discovered = true

			} else if !neighbor.Discovered && neighbor.Height < knot.Height {
				if *descender == nil || neighbor.Height > (*descender).Height {
					*descender = neighbor
					(*descender).Precursor = knot
					(*descender).Distance = knot.Distance + 1
					(*descender).Discovered = true
				}
			}

			if neighbor.Discovered && neighbor.Label == "end" {
				*end = neighbor
			}
		}
	}

	if len(*queue) > 0 {
		next := dequeue(queue)
		_DiscoverPlane(next, ascender, descender, end, queue)
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

func GetStart(net [][]*Knot) (x, y int) {
	for y := 0; y < len(net)-1; y++ {
		for x := 0; x < len((net)[0])-1; x++ {
			if (net)[y][x].Label == "start" {
				return x, y
			}
		}
	}
	return -1, -1
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
