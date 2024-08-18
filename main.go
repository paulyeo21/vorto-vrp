package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const MAX_DRIVE_TIME = 12 * 60

type Point struct {
	x, y float64
}

func (p1 *Point) Less(p2 *Point) bool {
	return distanceBetweenPoints(&Point{0, 0}, p1) < distanceBetweenPoints(&Point{0, 0}, p2)
}

func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

type Load struct {
	id              string
	pickup, dropoff *Point
}

func (l1 *Load) Less(l2 *Load) bool {
	return distanceBetweenPoints(&Point{0, 0}, l1.pickup) < distanceBetweenPoints(&Point{0, 0}, l2.pickup)
}

func (l Load) String() string {
	return fmt.Sprintf("%s -> Pickup: %v, Dropoff: %v", l.id, l.pickup, l.dropoff)
}

type Driver struct {
	position       *Point
	schedule       []string
	totalDriveTime float64
}

func newDriver() *Driver {
	return &Driver{
		position:       &Point{0, 0},
		schedule:       make([]string, 0),
		totalDriveTime: 0,
	}
}

func (d Driver) isNew() bool {
	return d.position.x == 0 && d.position.y == 0
}

func (d *Driver) completeLoad(l *Load) {
	d.totalDriveTime += distanceBetweenPoints(d.position, l.pickup)
	d.totalDriveTime += distanceBetweenPoints(l.pickup, l.dropoff)
	d.position = l.dropoff
	d.schedule = append(d.schedule, l.id)
}

func (d Driver) loadExceedsMaxDriveTime(load *Load) bool {
	distance := distanceBetweenPoints(d.position, load.pickup)
	distance += distanceBetweenPoints(load.pickup, load.dropoff)
	distance += distanceBetweenPoints(load.dropoff, &Point{0, 0})
	return d.totalDriveTime+distance > MAX_DRIVE_TIME
}

func (d *Driver) goHome() {
	d.totalDriveTime += distanceBetweenPoints(d.position, &Point{0, 0})
}

func (d Driver) printSchedule() {
	var buf bytes.Buffer
	buf.WriteByte('[')

	for i, id := range d.schedule {
		buf.WriteString(id)

		if i != len(d.schedule)-1 {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')
	buf.WriteByte('\n')
	buf.WriteTo(os.Stdout)
}

func (d Driver) String() string {
	return fmt.Sprintf("%v %v %f", d.position, d.schedule, d.totalDriveTime)
}

func distanceBetweenPoints(p1, p2 *Point) float64 {
	return math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
}

func getPointFromPointStr(s string) *Point {
	r := strings.NewReplacer(" ", "", "(", "", ")", "")
	split := strings.Split(r.Replace(s), ",")

	x, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		panic(err)
	}

	y, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		panic(err)
	}

	return &Point{x, y}
}

func assignDriversToLoads(loads []*Load) {
	// Build BST with loads sorted by pickup distance to origin (0, 0)
	bst := &Node{Load: loads[0]}

	for i, load := range loads {
		if i == 0 {
			continue
		}

		bst.Insert(&Node{Load: load})
	}

	// Search for load pickup closest to drivers current position
	driver := newDriver()

	for i := 0; i < len(loads); i++ {
		node := bst.Search(driver.position)
		// fmt.Println("Driver: ", driver)

		// if distance between (node and driver position + home to node pickup) > 500
		// allocate new driver

		if driver.loadExceedsMaxDriveTime(node.Load) {
			driver.goHome()
			driver.printSchedule()

			// Allocate new driver
			driver = newDriver()
			node = bst.Min()
		}

		// fmt.Println("Node: ", node)
		driver.completeLoad(node.Load)
		bst = bst.Delete(node)
		// fmt.Println(bst.Delete(node))
	}

	driver.goHome()
	driver.printSchedule()
}

/*
	i.e. go run main.go <path>
*/
func main() {
	// 1. Open problem file
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 2. Read loads line by line
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Assume first line is header and skip

	loads := make([]*Load, 0)

	for scanner.Scan() {
		load := strings.Split(scanner.Text(), " ")
		id, pickup, dropoff := load[0], getPointFromPointStr(load[1]), getPointFromPointStr(load[2])
		loads = append(loads, &Load{id, pickup, dropoff})
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// 3. Print out schedule
	assignDriversToLoads(loads)
}
