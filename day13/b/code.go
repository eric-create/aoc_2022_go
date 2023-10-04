package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")
	lines = append(lines, "[[2]]", "[[6]]")
	packets := SortPackets(lines)

	two := -1
	six := -1

	for i, packet := range *packets {
		fmt.Println(i, ListToString(packet))

		switch packetString := ListToString(packet); packetString {
		case "[[2]]":
			two = i + 1
		case "[[6]]":
			six = i + 1
		}
	}

	fmt.Println("two", two, "times", "six", six, "equals", two*six)
}

func SortPackets(lines []string) *[]NestedList {
	packets := GetPackets(lines)
	sortedPackets := []NestedList{packets[0]}
	packets = packets[1:]

	for _, packet := range packets {
		for i, sorted := range sortedPackets {
			if Compare(packet, sorted) == Success {
				sortedPackets = slices.Insert(sortedPackets, i, packet)
				break
			} else if i == len(sortedPackets)-1 {
				sortedPackets = append(sortedPackets, packet)
			}
		}
	}

	return &sortedPackets
}

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

	stepNumber := 0
	for cursor := p.position; cursor <= p.end; cursor++ {
		currentSymbol := string(p.Data[cursor])
		if matched, _ := regexp.MatchString("[0-9]", currentSymbol); matched {
			number += currentSymbol
			stepNumber++
		} else {
			break
		}
	}
	p.position += stepNumber - 1
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

type Result int

const (
	Success Result = iota
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
		return Success
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
		if i == len(rightList.elements) {
			return Error
		}

		result := Compare(leftElement, rightList.elements[i])

		// Only return a "breaking result", that is the left number being smaller or bigger
		// than the right number.
		// If the left number equals the right number, continue comparing the next elements
		// in the list.
		if result == Error || result == Success {
			return result
		}
	}

	// The left list ran out of elements and is of the same size as the right list.
	if len(leftList.elements) == len(rightList.elements) {
		return Continue
	}

	// The left list ran out of elements, but is smaller than the right list.
	return Success
}

func GetPackets(lines []string) []NestedList {
	packets := []NestedList{}

	for _, line := range lines {
		if line != "" {
			packet := NewPacketIterator(line).MakeNestedList()
			packets = append(packets, packet)
		}
	}

	return packets
}

func ListToString(nl NestedList) string {
	packetString := ""
	PacketString(nl, &packetString)
	packetString = packetString[:len(packetString)-1]
	return packetString
}

func PacketString(nl NestedList, packetString *string) {
	if nl.IsList() {
		*packetString += "["
		list := nl.List()
		for _, e := range list.elements {
			PacketString(e, packetString)
		}
		if len(list.elements) > 0 {
			*packetString = string((*packetString)[:len(*packetString)-1])
		}
		*packetString += "],"
	} else {
		integer := nl.Integer()
		*packetString += fmt.Sprintf("%d,", integer.Value)
	}
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
