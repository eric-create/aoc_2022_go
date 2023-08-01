package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Symbol struct {
	Value any
}

func NewSymbol(s string) *Symbol {
	symbol := Symbol{}

	if matched, _ := regexp.MatchString("[0-9]+", s); matched {
		symbol.Value = func() int { integer, _ := strconv.Atoi(s); return integer }()
	} else {
		symbol.Value = s
	}

	return &symbol
}

func (s *Symbol) isInteger() bool {
	return reflect.TypeOf((*s).Value).Kind() == reflect.Int
}

func (s *Symbol) Integer() *Integer {
	var i Integer
	if str, ok := s.Value.(int); ok {
		i = Integer{int(str)}
		return &i
	} else {
		log.Fatal("Severe error!")
		// The following code should actually never be reached. Don't know why the compiler
		// is stupid.
		i = Integer{-1}
		return &i
	}
}

func (s *Symbol) isOpenBracket() bool {
	if !s.isInteger() {
		return s.Value == "["
	}
	return false
}

func (s *Symbol) isCloseBracket() bool {
	if !s.isInteger() {
		return s.Value == "]"
	}
	return false
}

type PacketIterator struct {
	Data     string
	position int
	end      int
}

func NewPacketIterator(data string) *PacketIterator {
	p := PacketIterator{Data: data}
	// It is clear that the first position of the packet string is an opening bracket.
	// Therefore return that into the Nirvana.
	p.position = 0
	// It is clear that the last position in the packet string holds a closing bracket.
	p.end = len(p.Data) - 2

	return &p
}

// Return the single digit string at the current position of the iterator.
func (p *PacketIterator) CurrentCharacter() string {
	return string(p.Data[p.position])
}

// Advances the position of the cursor of the `PacketIterator` by one. Returns the symbol
// at that position. If the cursor has already reached the end of the packet, then returns
// nil.
func (p *PacketIterator) NextSymbol() *Symbol {

	// If the end of the packet is reached, return nil.
	if p.position == p.end {
		return nil
	}

	// Advance to the next position in the packet string.
	p.position++

	currentCharacter := p.CurrentCharacter()

	// Return brackets.
	if matched, _ := regexp.MatchString("[\\[\\]]", currentCharacter); matched {
		return NewSymbol(currentCharacter)
	}

	// Return number.
	if matched, _ := regexp.MatchString("[0-9]+", currentCharacter); matched {
		return p.CurrentNumber()
	}

	// Advance to the next position once again if a comma. It is clear at this point that
	// the current symbol must be a comma, since all existing symbols are brackets, numbers
	// and commas.
	return p.NextSymbol()
}

// Gets the number at the current position of the iterator its cursor.
func (p *PacketIterator) CurrentNumber() *Symbol {
	number := ""
	for cursor := p.position; cursor <= p.end; cursor++ {
		currentSymbol := string(p.Data[cursor])
		if matched, _ := regexp.MatchString("[0-9]", currentSymbol); matched {
			number += currentSymbol
		} else {
			break
		}
	}
	return NewSymbol(number)
}

func (p *PacketIterator) CurrentSymbol() *Symbol {
	// Return brackets.
	if matched, _ := regexp.MatchString("[\\[\\]]]", p.CurrentCharacter()); matched {
		return NewSymbol(p.CurrentCharacter())
	}

	return p.CurrentNumber()
}

func (p *PacketIterator) MakeNestedList() NestedList {
	initialList := List{}
	makeNestedList(p, &initialList)
	return &initialList
}

func makeNestedList(p *PacketIterator, list *List) {
	for symbol := p.NextSymbol(); symbol != nil; symbol = p.NextSymbol() {

		if symbol.isInteger() {
			list.AppendInteger(symbol.Integer())

		} else if symbol.isOpenBracket() {
			nestedList := List{}
			makeNestedList(p, &nestedList)
			list.AppendList(&nestedList)

		} else if symbol.isCloseBracket() {
			return
		}
	}
}

type NestedList interface {
	IsList() bool
	IsInteger() bool
	List() *List
	Integer() *Integer
}

type Integer struct {
	Value int
}

func (i *Integer) IsList() bool {
	return false
}

func (i *Integer) IsInteger() bool {
	return true
}

func (i *Integer) List() *List {
	panic("This is of type Integer, not of type List!")
}

func (i *Integer) Integer() *Integer {
	return i
}

type List struct {
	elements []NestedList
}

func (l *List) IsList() bool {
	return true
}

func (l *List) IsInteger() bool {
	return false
}

func (l *List) List() *List {
	return l
}

func (l *List) Integer() *Integer {
	panic("This is of type List, not of type Integer!")
}

func (l *List) AppendInteger(i *Integer) {
	l.elements = append(l.elements, i)
}

func (l *List) AppendList(list *List) {
	l.elements = append(l.elements, list)
}

func main() {
	lines := ReadLines("./input.txt")
	listPairs := getListPairs(lines)
	sum := 0

	for i, listPair := range listPairs {
		index := i + 1
		isCorrectlyOrdered := Compare(listPair[0], listPair[1]) == Sucess
		if isCorrectlyOrdered {
			fmt.Println(index, isCorrectlyOrdered)
			sum += index
		}
	}
	fmt.Println("\n", sum)
}

type Result int

const (
	Sucess Result = iota
	Continue
	Error
)

func Compare(left NestedList, right NestedList) Result {

	if left.IsInteger() && right.IsInteger() { // Both elements are integers.
		return CompareIntegers(left, right)

	} else if left.IsList() && right.IsList() { // Both elements are lists.
		return CompareLists(left, right)

	} else { // One element is a list and one is an integer.

		// Left element is a list and right element is an integer.
		if left.IsList() && right.IsInteger() {
			newRightList := List{[]NestedList{right.Integer()}}
			return CompareLists(left, &newRightList)

		} else { // Left element is an integer and right element is a list.
			newLeftList := List{[]NestedList{left.Integer()}}
			return CompareLists(&newLeftList, right)
		}
	}
}

func CompareIntegers(left NestedList, right NestedList) Result {
	leftInt := left.Integer().Value
	rightInt := right.Integer().Value

	if leftInt < rightInt {
		return Sucess
	} else if leftInt == rightInt {
		return Continue
	} else {
		return Error
	}
}

func CompareLists(left NestedList, right NestedList) Result {
	leftList := left.List()
	rightList := right.List()

	for i, leftElement := range leftList.elements {

		// If the right list runs out of elements, this is an Error.
		if i >= len(rightList.elements) {
			return Error
		}

		result := Compare(leftElement, rightList.elements[i])

		// Only return a "breaking result", that is the left number being smaller or bigger
		// than the right number.
		// If the left number equals the right number, continue comparing the next elements
		// in the list.
		if result == Error || result == Sucess {
			return result
		}

		// If all the element pairs are equal and both lists run out of elements at the
		// same time, that is Continue.
		if i == len(leftList.elements)-1 && i == len(rightList.elements)-1 {
			return Continue
		}
	}

	// The left list ran out of elements, this is a Success.
	return Sucess
}

func getListPairs(lines []string) [][2]NestedList {
	listPairs := [][2]NestedList{}
	pairs := getStringPairs(lines)

	for _, pair := range pairs {
		left := NewPacketIterator(pair[0]).MakeNestedList()
		PrintList(left)
		fmt.Println()
		right := NewPacketIterator(pair[1]).MakeNestedList()
		PrintList(right)
		fmt.Println()
		fmt.Println()
		listPairs = append(listPairs, [2]NestedList{left, right})
	}

	return listPairs
}

func PrintList(nl NestedList) {
	if nl.IsList() {
		fmt.Print("[")
		list := nl.List()
		for _, e := range list.elements {
			PrintList(e)
		}
		fmt.Print("]")
	} else {
		integer := nl.Integer()
		fmt.Print(integer.Value, ",")
	}
}

func getStringPairs(lines []string) [][2]string {

	lines = removeBlankLines(lines)

	pairs := [][2]string{}
	pair := [2]string{}

	for i, line := range lines {
		pair[i%2] = strings.ReplaceAll(line, " ", "")

		if i%2 == 1 {
			pairs = append(pairs, pair)
			pair = [2]string{}
		}
	}

	return pairs
}

func removeBlankLines(lines []string) []string {
	cleansed := []string{}

	for _, line := range lines {
		if line != "" {
			cleansed = append(cleansed, line)
		}
	}

	return cleansed
}

func ReadLines(path string) []string {
	content, _ := os.ReadFile(path)
	return strings.Split(string(content), "\n")
}
