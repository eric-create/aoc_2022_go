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

	node := NewNode("/", nil)
	Populate(&node, &lines)
	UpdateSizes(&node)
	// PrintTree(node, "")
	fmt.Println("---")
	fmt.Println(GetCandidates(node, 0))
}

func GetCandidates(node Node, sum int) int {
	for _, node := range node.Children {
		sum += GetCandidates(*node, 0)
	}
	if node.Size < 100000 {
		return sum + node.Size
	} else {
		return sum
	}
}

func PrintTree(node Node, offset string) {
	fmt.Println(offset+node.Name, "( size:", node.Size, ")")
	for filename, filesize := range node.Files {
		fmt.Println(offset+"  "+filename, filesize)
	}
	for _, child := range node.Children {
		PrintTree(*child, offset+"  ")
	}
}

func UpdateSizes(node *Node) {
	sum := 0
	for _, child := range (*node).Children {
		UpdateSizes(child)
		sum += (*child).Size
	}
	(*node).Size += sum
}

// Pops the element at index 0.
func PopLine(lines *[]string) string {
	if len(*lines) > 0 {
		line := (*lines)[0]
		*lines = (*lines)[1:]
		return line
	}
	return ""
}

// Returns the element at index 1.
func NextLine(lines *[]string) string {
	if len(*lines) > 0 {
		return (*lines)[0]
	}
	return ""
}

func isCd(line string) bool {
	return strings.HasPrefix(line, "$ cd")
}

func isLs(line string) bool {
	return line == "$ ls"
}

func Populate(node *Node, lines *[]string) {
	if line := PopLine(lines); line != "" {

		if isCd(line) {
			node = cd(node, line)

		} else if isLs(line) {
			AppendElements(node, lines)
		}

		Populate(node, lines)
	}
}

func AppendElements(node *Node, lines *[]string) {
	nextLine := NextLine(lines)
	for nextLine != "" && !isCd(nextLine) && !isLs(nextLine) {
		AppendElement(node, lines)
		nextLine = NextLine(lines)
	}
}

func AppendElement(node *Node, lines *[]string) {
	line := PopLine(lines)

	if strings.HasPrefix(line, "dir ") {
		newChild := NewNode(line[4:], node)
		(*node).Children = append((*node).Children, &newChild)

	} else {
		filesize, _ := strconv.Atoi(strings.Split(line, " ")[0])
		filename := strings.Split(line, " ")[1]
		(*node).Files[filename] = filesize
		(*node).Size += filesize
	}
}

func cd(node *Node, line string) *Node {
	destination := line[5:]

	switch destination {
	case "..":
		return (*node).Parent
	case "/":
		return node.root()
	}
	return (*node).get_node(destination)
}

type Node struct {
	Name     string
	Parent   *Node
	Children []*Node
	Files    map[string]int
	Size     int
}

func (node *Node) root() *Node {
	if node.Parent == nil {
		return node
	}
	return node.Parent.root()
}

func (node *Node) get_node(name string) *Node {
	for _, child := range (*node).Children {
		if child.Name == name {
			return child
		}
	}
	return node
}

func NewNode(name string, parent *Node) Node {
	newNode := Node{
		Name:     name,
		Parent:   parent,
		Children: []*Node{},
		Files:    map[string]int{},
		Size:     0}
	return newNode
}

func ReadFile(path string) []string {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(file), "\n")
}
