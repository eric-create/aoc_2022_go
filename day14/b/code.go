package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("input.txt")
	paths := GetPaths(lines)
	xMax, yMax := GetFieldConstraints(paths)
	field := InitField(paths, xMax, yMax)

	SetRocks(paths, &field)
	// grainSources := [][2]int{{500, 0}}
	sum := DropGrains(&field, [2]int{500, 0})
	PrintField(&field)

	fmt.Println()
	fmt.Println("number of units:", sum)
}

func DropGrains(field *[][]string, source [2]int) int {
	xStart := source[0]
	yStart := source[1]
	stopX, stopY := -1, -1
	sum := 0

	for stopX != 500 || stopY != 0 {
		// fmt.Println()
		// PrintField(field)

		stopX, stopY = DropGrain(field, xStart, yStart)
		sum++
	}

	return sum
}

func DropGrain(field *[][]string, x, y int) (int, int) {
	nextX, nextY := NextPosition(field, x, y)

	for nextX != x || nextY != y {
		x = nextX
		y = nextY
		nextX, nextY = NextPosition(field, nextX, nextY)
	}

	(*field)[y][x] = "o"

	return x, y
}

func NextPosition(field *[][]string, x, y int) (int, int) {
	// Reached the bottom
	if y+1 == len(*field) {
		return x, y
	}
	// Down
	if (*field)[y+1][x] == "." {
		return x, y + 1
	}
	// Down Left
	if (*field)[y+1][x-1] == "." {
		return x - 1, y + 1
	}
	// Down Right
	if (*field)[y+1][x+1] == "." {
		return x + 1, y + 1
	}
	// Rest
	return x, y
}

func PrintField(field *[][]string) {
	for y := 0; y < len(*field); y++ {
		for x := 480; x <= 520; x++ {
			fmt.Print((*field)[y][x])
		}
		fmt.Println()
	}
}

func SetRocks(paths []*Path, field *[][]string) {
	for _, path := range paths {
		for _, coordinate := range path.Coordinates {
			(*field)[coordinate[1]][coordinate[0]] = "#"
		}
	}
}

func InitField(paths []*Path, xMax, yMax int) [][]string {

	field := [][]string{}

	for y := 0; y <= yMax; y++ {
		field = append(field, []string{})

		for x := 0; x <= xMax; x++ {
			if y == yMax {
				field[y] = append(field[y], "#")
			} else {
				field[y] = append(field[y], ".")
			}

		}
	}

	return field
}

func GetFieldConstraints(paths []*Path) (int, int) {
	x, y := 1000, 0
	for _, path := range paths {
		for _, coordinate := range path.Coordinates {
			if coordinate[1] > y {
				y = coordinate[1]
			}
		}
	}

	return x, y + 2
}

func GetPaths(lines []string) []*Path {
	paths := []*Path{}

	for _, line := range lines {
		paths = append(paths, GetPath(line))
	}

	return paths
}

func GetPath(line string) *Path {
	path := Path{}

	coordinates := strings.Split(line, " -> ")
	for _, coordinate := range coordinates {
		components := strings.Split(coordinate, ",")
		x, _ := strconv.Atoi(components[0])
		y, _ := strconv.Atoi(components[1])
		(&path).Extend(x, y)
	}

	return &path
}

type Path struct {
	Coordinates [][2]int
}

func (path *Path) Last() (int, int) {
	if len((*path).Coordinates) == 0 {
		log.Fatal("called Path.Last() but there are no coordinates in Path")
	}

	last := (*path).Coordinates[len((*path).Coordinates)-1]
	return last[0], last[1]
}

func (path *Path) Extend(x, y int) {
	if len((*path).Coordinates) == 0 {
		(*path).Coordinates = append((*path).Coordinates, [2]int{x, y})

	} else {
		xStart, yStart := path.Last()
		xDir, yDir := Direction(xStart, yStart, x, y)

		for lastX, lastY := path.Last(); lastX != x || lastY != y; lastX, lastY = path.Last() {

			newCoordinate := [2]int{lastX + xDir, lastY + yDir}
			(*path).Coordinates = append((*path).Coordinates, newCoordinate)
		}
	}
}

func Direction(xStart, yStart, xEnd, yEnd int) (int, int) {
	xDir := NormalizeInt(xEnd - xStart)
	yDir := NormalizeInt(yEnd - yStart)

	return xDir, yDir
}

func NormalizeInt(i int) int {
	if i < 0 {
		return -1
	} else if i > 0 {
		return 1
	}
	return 0
}

func ReadLines(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("failed to read ", path)
	}

	return strings.Split(string(content), "\n")
}
