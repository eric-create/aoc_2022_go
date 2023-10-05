package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	lines := ReadLines("./input.txt")

	sensors := GetSensors(lines)
	xMin, yMin, xMax, yMax := EventHorizon(sensors)
	numCovered := Coverage(sensors, xMin, xMax, 2000000) // 2000000 or 10
	fmt.Println(numCovered)
	fmt.Println(xMin, yMin, xMax, yMax)
}

func Coverage(sensors *[]*Sensor, xMin, xMax, y int) int {
	coverage := []int{}

	for x := xMin; x <= xMax; x++ {
		for _, sensor := range *sensors {
			if sensor.Covers(Coordinate{x, y}) {
				if len(coverage) == 0 || coverage[len(coverage)-1] != x {
					coverage = append(coverage, x)
				}
			}
		}
	}

	count := len(coverage)
	beacons := []*Coordinate{}

	for _, sensor := range *sensors {
		beacon := sensor.Beacon

		if beacon.Y == y && beacon.X >= xMin && beacon.X <= xMax {
			if !slices.Contains(beacons, beacon) {
				beacons = append(beacons, beacon)
				count--
			}
		}
	}

	return count
}

func EventHorizon(sensors *[]*Sensor) (int, int, int, int) {
	xMin, yMin := (*sensors)[0].Position.X, (*sensors)[0].Position.Y
	xMax, yMax := 0, 0

	for _, sensor := range *sensors {
		position := sensor.Position

		if position.X+sensor.Radius > xMax {
			xMax = position.X + sensor.Radius
		}

		if position.Y+sensor.Radius > yMax {
			yMax = position.Y + sensor.Radius
		}

		if position.X-sensor.Radius < xMin {
			xMin = position.X - sensor.Radius
		}

		if position.Y-sensor.Radius < yMin {
			yMin = position.Y - sensor.Radius
		}
	}

	return xMin, yMin, xMax, yMax
}

func GetSensors(lines []string) *[]*Sensor {
	sensors := []*Sensor{}
	beacons := []*Coordinate{}

	for _, line := range lines {
		sensor := GetSensor(line, &beacons)
		sensors = append(sensors, sensor)
	}

	return &sensors
}

func GetSensor(line string, beacons *[]*Coordinate) *Sensor {
	components := strings.Split(line, ":")
	prefix := components[0]
	suffix := components[1]

	sensorCoordinate := NewCoordinate(prefix[10:])
	beaconCoordinate := NewCoordinate(suffix[22:])

	beacon := GetBeacon(&beaconCoordinate, beacons)

	sensor := NewSensor(sensorCoordinate, beacon)

	return sensor
}

func GetBeacon(newBeacon *Coordinate, beacons *[]*Coordinate) *Coordinate {
	for _, beacon := range *beacons {
		if beacon.X == newBeacon.X && beacon.Y == newBeacon.Y {
			return beacon
		}
	}

	*beacons = append(*beacons, newBeacon)
	return newBeacon
}

func NewCoordinate(s string) Coordinate {
	regex := regexp.MustCompile(`x=(-?\d+), y=(-?\d+)`)
	matches := regex.FindStringSubmatch(s)

	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])

	return Coordinate{x, y}
}

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) Navigate(direction Coordinate, times int) *Coordinate {
	new := Coordinate{c.X, c.Y}

	for i := 0; i < times; i++ {
		new.X += direction.X
		new.Y += direction.Y
	}

	return &new
}

type Sensor struct {
	Position Coordinate
	Beacon   *Coordinate
	Radius   int
}

func (s *Sensor) Covers(position Coordinate) bool {
	distance := Distance(s.Position, position)
	radius := s.Radius
	return distance <= radius
}

func NewSensor(position Coordinate, beacon *Coordinate) *Sensor {
	return &Sensor{
		Position: position,
		Beacon:   beacon,
		Radius:   Distance(position, *beacon)}
}

func Distance(start, end Coordinate) int {
	xDiff := int(math.Abs(float64(start.X) - float64(end.X)))
	yDiff := int(math.Abs(float64(start.Y) - float64(end.Y)))
	distance := xDiff + yDiff
	return distance
}

func ReadLines(path string) []string {
	content, _ := os.ReadFile(path)
	return strings.Split(string(content), "\n")
}
