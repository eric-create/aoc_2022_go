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
	SetScenicScores(grid)
	x, y := GetHigh(grid)
	fmt.Println(x, y, grid[y][x].Score)
}

func GetHigh(grid [][]*Tree) (int, int) {
	highTree := grid[0][0]
	xPos := 0
	yPos := 0

	for y := range grid {
		for x := range grid[y] {
			if grid[y][x].Score > highTree.Score {
				highTree = grid[y][x]
				yPos = y
				xPos = x
			}
		}
	}

	return xPos, yPos
}

func SetScenicScores(grid [][]*Tree) {
	for y := range grid {
		for x := range grid[y] {
			upScore := GetScenicScore(grid, grid[y][x], 0, -1, x, y)
			downScore := GetScenicScore(grid, grid[y][x], 0, 1, x, y)
			leftScore := GetScenicScore(grid, grid[y][x], -1, 0, x, y)
			rightScore := GetScenicScore(grid, grid[y][x], 1, 0, x, y)
			grid[y][x].Score = upScore * downScore * leftScore * rightScore
		}
	}
}

func GetScenicScore(grid [][]*Tree, tree *Tree, xDir int, yDir int, xPos int, yPos int) int {
	yEdge := len(grid) - 1
	xEdge := len(grid[0]) - 1

	if isEdge(xEdge, yEdge, xPos, yPos) {
		return 0
	}
	if isEdge(xEdge, yEdge, xPos+xDir, yPos+yDir) {
		return 1
	}
	if grid[yPos+yDir][xPos+xDir].Height >= tree.Height {
		return 1
	}
	return 1 + GetScenicScore(grid, tree, xDir, yDir, xPos+xDir, yPos+yDir)
}

func isEdge(xEdge int, yEdge int, xPos int, yPos int) bool {
	if xPos == 0 || xPos == xEdge || yPos == 0 || yPos == yEdge {
		return true
	}
	return false
}

type Tree struct {
	Height int
	Score  int
}

func Grid(lines []string) [][]*Tree {
	grid := [][]*Tree{}
	for i, line := range lines {
		grid = append(grid, []*Tree{})

		for _, char := range line {
			height, _ := strconv.Atoi(string(char))
			grid[i] = append(grid[i], &Tree{height, 0})
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
