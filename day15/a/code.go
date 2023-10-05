package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")

	sensors := GetSensors(lines)

	xMin, yMin, xMax, yMax := EventHorizon(sensors)
	field := InitField(xMin, yMin, xMax, yMax)
	SetField(sensors, field)

	for _, sensor := range *sensors {
		sensor.Sense(field)
	}

	numCovered := NumCovered(10, field)
	fmt.Println(numCovered)

	// PrintField(field)
}

func EventHorizon(sensors *[]*Sensor) (int, int, int, int) {
	xMin, yMin := (*sensors)[0].Position.X, (*sensors)[0].Position.Y
	xMax, yMax := 0, 0

	for _, sensor := range *sensors {
		pos := sensor.Position

		if pos.X > xMax {
			xMax = pos.X
		}

		if pos.Y > yMax {
			yMax = pos.Y
		}

		if pos.X < xMin {
			xMin = pos.X
		}

		if pos.Y < yMin {
			yMin = pos.Y
		}
	}

	return xMin - 9, yMin - 9, xMax + 9, yMax + 9
}

func NumCovered(y int, field *[][]*Coordinate) int {
	xMax := len((*field)[0]) - 1
	covered := 0

	for x := 0; x <= xMax; x++ {
		symbol := (*field)[Norm(y)][x].Symbol

		if symbol == "#" {
			covered++
		}
	}

	return covered
}

func SetField(sensors *[]*Sensor, field *[][]*Coordinate) {
	for _, sensor := range *sensors {

		sensorPos := &sensor.Position
		(*field)[sensorPos.NormY()][sensorPos.NormX()] = sensorPos

		beaconPos := sensor.Beacons[0]
		(*field)[beaconPos.NormY()][beaconPos.NormX()] = beaconPos
	}
}

func PrintField(field *[][]*Coordinate) {
	for y := 0; y < len(*field); y++ {
		for x := 0; x < len((*field)[0]); x++ {
			fmt.Print((*field)[y][x].Symbol)
		}
		fmt.Println()
	}
}

func InitField(xMin, yMin, xMax, yMax int) *[][]*Coordinate {
	field := [][]*Coordinate{}

	for y := yMin; y <= yMax; y++ {
		field = append(field, []*Coordinate{})
		for x := xMin; x <= xMax; x++ {
			if x == 0 && y == 0 {
				field[Norm(y)] = append(field[Norm(y)], &Coordinate{0, 0, "M"})
			} else {
				field[Norm(y)] = append(field[Norm(y)], &Coordinate{x, y, "."})
			}
		}
	}

	return &field
}

func GetFieldConstraints(sensors *[]*Sensor) (int, int, int, int) {
	xMin, yMin, xMax, yMax := 0, 0, 0, 0

	for _, sensor := range *sensors {
		pos := (*sensor).Position
		if pos.X > xMax {
			xMax = pos.X
		}
		if pos.Y > yMax {
			yMax = pos.Y
		}
		if pos.X < xMin {
			xMin = pos.X
		}
		if pos.Y < yMin {
			yMin = pos.Y
		}
	}

	// The sensor can see a distance of 9.
	return xMin - 9, yMin - 9, xMax + 9, yMax + 9
}

func GetSensors(lines []string) *[]*Sensor {
	sensors := []*Sensor{}

	for _, line := range lines {
		sensor := GetSensor(line)
		sensors = append(sensors, sensor)
	}

	return &sensors
}

func GetSensor(line string) *Sensor {
	components := strings.Split(line, ":")
	prefix := components[0]
	suffix := components[1]

	sensorCoordinate := GetCoordinate(prefix[10:], "S")
	beaconCoordinate := GetCoordinate(suffix[22:], "B")

	sensor := Sensor{
		Position: sensorCoordinate,
		Beacons:  []*Coordinate{&beaconCoordinate},
	}

	return &sensor
}

func GetCoordinate(s string, symbol string) Coordinate {
	regex := regexp.MustCompile(`x=(-?\d+), y=(-?\d+)`)
	matches := regex.FindStringSubmatch(s)

	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])

	return Coordinate{x, y, symbol}
}

type Coordinate struct {
	X      int
	Y      int
	Symbol string
}

func (c *Coordinate) Navigate(direction Coordinate, times int) *Coordinate {
	new := Coordinate{c.X, c.Y, ""}

	for i := 0; i < times; i++ {
		new.X += direction.X
		new.Y += direction.Y
	}

	return &new
}

func (c *Coordinate) NormX() int {
	return Norm((*c).X)
}

func (c *Coordinate) NormY() int {
	return Norm((*c).Y)
}

func Norm(i int) int {
	return i + 9
}

type Sensor struct {
	Position Coordinate
	Beacons  []*Coordinate
}

func (s *Sensor) Sense(field *[][]*Coordinate) {
	movements := [][2]Coordinate{
		// Right-Up
		{Coordinate{1, 0, ""}, Coordinate{0, -1, ""}},
		// Down-Right
		{Coordinate{0, 1, ""}, Coordinate{1, 0, ""}},
		// Left-Down
		{Coordinate{-1, 0, ""}, Coordinate{0, 1, ""}},
		// Up-Left
		{Coordinate{0, -1, ""}, Coordinate{-1, 0, ""}}}

	detected := false

	for i := 0; i <= 9; i++ {

		for _, movement := range movements {
			priMovement := movement[0]
			secMovement := movement[1]

			for j := 0; j <= i; j++ {
				targetPos := s.Position.Navigate(priMovement, j).Navigate(secMovement, i-j)
				cursor := (*field)[targetPos.NormY()][targetPos.NormX()]

				if cursor.Symbol == "." {
					cursor.Symbol = "#"

				} else if cursor.Symbol == "B" {
					detected = true

				}
			}
		}

		if detected {
			return
		}
	}
}

func ReadLines(path string) []string {
	content, _ := os.ReadFile(path)
	return strings.Split(string(content), "\n")
}
