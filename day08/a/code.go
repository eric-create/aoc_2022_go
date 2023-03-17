package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := ReadFile("./input.txt")
	grid := Grid(lines)

	allVisibles := []*Tree{}

	for i := 0; i < 4; i++ {
		allVisibles = extend(allVisibles, VisibleTrees(grid))
		grid = Rotate(grid)
	}

	fmt.Println(len(allVisibles))
	// fmt.Println()
	// PrintGrid(grid, allVisibles)
}

func PrintGrid(grid [][]*Tree, visibleTrees []*Tree) {
	for i := range grid {
		for j := range grid[i] {
			if contains(visibleTrees, grid[i][j]) {
				fmt.Print(grid[i][j].Height, " ")
			} else {
				fmt.Print(" ", " ")
			}
		}
		fmt.Print("\n")
	}
}

func MakeGrid(y_size int, x_size int) [][]*Tree {
	grid := [][]*Tree{}
	for y := 0; y < y_size; y++ {
		grid = append(grid, make([]*Tree, x_size))
	}
	return grid
}

func Rotate(old [][]*Tree) [][]*Tree {
	// Create a new grid where the the x_size is the old y_size and vice versa.
	new := MakeGrid(len(old), len(old[0]))
	y_len := len(new)

	for y := range old {
		for x := range old[y] {
			new[x][y_len-1-y] = old[y][x]
		}
	}

	return new
}

func contains(trees []*Tree, tree *Tree) bool {
	for _, _tree := range trees {
		if _tree == tree {
			return true
		}
	}
	return false
}

func extend(allVisibles []*Tree, visibles []*Tree) []*Tree {
	for _, visible := range visibles {
		if !contains(allVisibles, visible) {
			allVisibles = append(allVisibles, visible)
		}
	}
	return allVisibles
}

func VisibleTrees(old [][]*Tree) []*Tree {
	allVisibles := []*Tree{}

	for _, row := range old {
		rowVisibles := RowVisibles(row)
		allVisibles = extend(allVisibles, rowVisibles)
	}

	// fmt.Println()
	return allVisibles
}

func RowVisibles(row []*Tree) []*Tree {
	visibles := []*Tree{}
	visibles = append(visibles, row[0])

	// fmt.Print((*visibles[0]).Height)
	for i, tree := range row[1:] {
		if tree.Height > visibles[len(visibles)-1].Height {
			// fmt.Print(tree.Height)
			visibles = append(visibles, row[i+1])
		} else {
			// fmt.Print(" ")
		}
	}
	// fmt.Print("\n")
	return visibles
}

type Tree struct {
	Height int
}

func Grid(lines []string) [][]*Tree {
	grid := [][]*Tree{}
	for i, line := range lines {
		grid = append(grid, []*Tree{})

		for _, char := range line {
			height, _ := strconv.Atoi(string(char))
			grid[i] = append(grid[i], &Tree{height})
		}
	}
	return grid
}

func ReadFile(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
