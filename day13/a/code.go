package main

import (
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
	if p.position+1 == p.end {
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
	isList() bool
	isInt() bool
}

type Integer struct {
	Value int
}

func (i *Integer) isList() bool {
	return false
}

func (i *Integer) isInt() bool {
	return true
}

type List struct {
	elements []NestedList
}

func (l *List) isList() bool {
	return true
}

func (l *List) isInt() bool {
	return false
}

func (l *List) AppendInteger(i *Integer) {
	l.elements = append(l.elements, i)
}

func (l *List) AppendList(nestedList *List) {
	l.elements = append(l.elements, nestedList)
}

func main() {
	lines := ReadLines("./input.txt")
	listPairs := getListPairs(lines)

	for _, listPair := range listPairs {
		isOrderedCorrectly(listPair[0], listPair[1])
	}
}

func getListPairs(lines []string) [][2]NestedList {
	listPairs := [][2]NestedList{}
	pairs := getStringPairs(lines)

	for _, pair := range pairs {
		left := NewPacketIterator(pair[0]).MakeNestedList()
		right := NewPacketIterator(pair[1]).MakeNestedList()
		listPairs = append(listPairs, [2]NestedList{left, right})
	}

	return listPairs
}

func isOrderedCorrectly(left NestedList, right NestedList) bool {
	return false
}

// func compare(left *List, right *List) bool {
// 	for leftElement := range left.elements {

// 	}
// }

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
