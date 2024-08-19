package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const MAX_DRIVE_TIME = 12 * 60 // 12 hours
const MIN_DETOUR_COST = 123    // tune

type Point struct {
	x, y float64
}

func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

type Load struct {
	id              string
	pickup, dropoff Point
}

func (l Load) distanceToPickup() float64 {
	return distanceBetweenPoints(Point{0.0, 0.0}, l.pickup)
}

func (l Load) distanceFromPickupToDropoff() float64 {
	return distanceBetweenPoints(l.pickup, l.dropoff)
}

func (l Load) distanceFromDropoffToHome() float64 {
	return distanceBetweenPoints(l.dropoff, Point{0.0, 0.0})
}

func (l Load) distanceToHome() float64 {
	return l.distanceToPickup() + l.distanceFromPickupToDropoff() + l.distanceFromDropoffToHome()
}

func (l Load) distanceFromHomeToDropoff() float64 {
	return l.distanceToPickup() + l.distanceFromPickupToDropoff()
}

func (l Load) distanceFromPickupToHome() float64 {
	return l.distanceFromPickupToDropoff() + l.distanceFromDropoffToHome()
}

func distanceBetweenPoints(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
}

func getPointFromPointStr(s string) Point {
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

	return Point{x, y}
}

func printSchedule(loads []Load) {
	var buf bytes.Buffer
	buf.WriteByte('[')

	for i, load := range loads {
		buf.WriteString(load.id)

		if i != len(loads)-1 {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')
	buf.WriteByte('\n')
	buf.WriteTo(os.Stdout)
}

func getDistanceOfScheduleWithReturnHome(schedule []Load) float64 {
	if len(schedule) == 0 {
		return 0.0
	}

	current := Point{0.0, 0.0}
	total := 0.0

	for _, load := range schedule {
		total += distanceBetweenPoints(current, load.pickup)
		total += load.distanceFromPickupToDropoff()
		current = load.dropoff
	}

	total += distanceBetweenPoints(current, Point{0.0, 0.0})
	return total
}

func assignDriversToLoads(loads []Load) {
	visited := make(map[string]struct{}, 0)

	sort.Slice(loads, func(i, j int) bool {
		return loads[i].distanceToHome() > loads[j].distanceToHome()
	})

	for i := 0; i < len(loads); i++ {
		if _, ok := visited[loads[i].id]; ok {
			continue
		}

		schedule := []Load{loads[i]}
		visited[loads[i].id] = struct{}{}

		// Find detours is minimized
		detourCost := math.MaxFloat64
		j := i

		// Find best detours from home to last load
		for j < len(loads) {
			for k := j + 1; k < len(loads); k += 1 {
				if _, ok := visited[loads[k].id]; ok {
					continue
				}

				// Check if adding load extends current schedule run time > 720
				if getDistanceOfScheduleWithReturnHome(append([]Load{loads[k]}, schedule...)) > MAX_DRIVE_TIME {
					continue
				}

				// Get cost of going to new load vs straight to previous load
				newDetourCost := loads[k].distanceFromHomeToDropoff() + distanceBetweenPoints(loads[k].dropoff, schedule[len(schedule)-1].pickup) - schedule[len(schedule)-1].distanceToPickup()

				if newDetourCost < detourCost {
					j = k
					detourCost = newDetourCost
				}

				if newDetourCost < MIN_DETOUR_COST {
					j = k
					break
				}
			}

			// Finish if no new detours
			if loads[j] == schedule[0] {
				break
			}

			schedule = append([]Load{loads[j]}, schedule...)
			visited[loads[j].id] = struct{}{}
		}

		detourCost = math.MaxFloat64
		j = i

		// Find best detours from last load to home
		for j < len(loads) {
			for k := j + 1; k < len(loads); k += 1 {
				if _, ok := visited[loads[k].id]; ok {
					continue
				}

				// Check if adding load extends current schedule run time > 720
				if getDistanceOfScheduleWithReturnHome(append(schedule, loads[k])) > MAX_DRIVE_TIME {
					continue
				}

				// Get cost of going to new load vs straight home
				newDetourCost := loads[k].distanceFromPickupToHome() + distanceBetweenPoints(schedule[len(schedule)-1].dropoff, loads[k].pickup) - schedule[len(schedule)-1].distanceFromDropoffToHome()

				if newDetourCost < detourCost {
					j = k
					detourCost = newDetourCost
				}

				if newDetourCost < MIN_DETOUR_COST {
					j = k
					break
				}
			}

			// Finish if no new detours
			if loads[j] == schedule[len(schedule)-1] {
				break
			}

			schedule = append(schedule, loads[j])
			visited[loads[j].id] = struct{}{}
		}

		printSchedule(schedule)
	}
}

func runStatistics(loads []Load) {
	sum := 0.0
	minimumDistance, maximumDistance := math.MaxFloat64, 0.0

	sumDistanceFromHomeToPickup := 0.0
	minimumDistanceFromHomeToPickup, maximumDistanceFromHomeToPickup := math.MaxFloat64, 0.0

	sumDistanceDropoffToHome := 0.0
	minimumDistanceDropoffToHome, maximumDistanceDropoffToHome := math.MaxFloat64, 0.0

	sumDistanceToHome := 0.0
	minimumDistanceToHome, maximumDistanceToHome := math.MaxFloat64, 0.0

	for _, load := range loads {
		distance := load.distanceFromPickupToDropoff()

		if distance < minimumDistance {
			minimumDistance = distance
		}

		if distance > maximumDistance {
			maximumDistance = distance
		}

		distanceFromHomeToPickup := load.distanceToPickup()

		if distanceFromHomeToPickup < minimumDistanceFromHomeToPickup {
			minimumDistanceFromHomeToPickup = distanceFromHomeToPickup
		}

		if distanceFromHomeToPickup > maximumDistanceFromHomeToPickup {
			maximumDistanceFromHomeToPickup = distanceFromHomeToPickup
		}

		distanceDropoffToHome := load.distanceFromDropoffToHome()

		if distanceDropoffToHome < minimumDistanceDropoffToHome {
			minimumDistanceDropoffToHome = distanceDropoffToHome
		}

		if distanceDropoffToHome > maximumDistanceDropoffToHome {
			maximumDistanceDropoffToHome = distanceDropoffToHome
		}

		distanceToHome := load.distanceToHome()

		if distanceToHome < minimumDistanceToHome {
			minimumDistanceToHome = distanceToHome
		}

		if distanceToHome > maximumDistanceToHome {
			maximumDistanceToHome = distanceToHome
		}

		sum += distance
		sumDistanceFromHomeToPickup += distanceFromHomeToPickup
		sumDistanceDropoffToHome += distanceDropoffToHome
		sumDistanceToHome += distanceToHome
	}

	fmt.Printf("Number of loads: %d\n", len(loads))

	fmt.Printf("Average distance of pickup to dropoff: %f\n", sum/float64(len(loads)))
	fmt.Printf("Minimum distance of pickup to dropoff: %f\n", minimumDistance)
	fmt.Printf("Maximum distance of pickup to dropoff: %f\n", maximumDistance)

	fmt.Printf("Average distance of pickup from home: %f\n", sumDistanceFromHomeToPickup/float64(len(loads)))
	fmt.Printf("Minimum distance of pickup from home: %f\n", minimumDistanceFromHomeToPickup)
	fmt.Printf("Maximum distance of pickup from home: %f\n", maximumDistanceFromHomeToPickup)

	fmt.Printf("Average distance of dropoff to home: %f\n", sumDistanceDropoffToHome/float64(len(loads)))
	fmt.Printf("Minimum distance of dropoff to home: %f\n", minimumDistanceDropoffToHome)
	fmt.Printf("Maximum distance of dropoff to home: %f\n", maximumDistanceDropoffToHome)

	fmt.Printf("Average distance to home: %f\n", sumDistanceToHome/float64(len(loads)))
	fmt.Printf("Minimum distance to home: %f\n", minimumDistanceToHome)
	fmt.Printf("Maximum distance to home: %f\n", maximumDistanceToHome)
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

	loads := make([]Load, 0)

	for scanner.Scan() {
		load := strings.Split(scanner.Text(), " ")
		id, pickup, dropoff := load[0], getPointFromPointStr(load[1]), getPointFromPointStr(load[2])
		loads = append(loads, Load{id, pickup, dropoff})
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// 3. Print out schedule
	// runStatistics(loads)
	assignDriversToLoads(loads)
}
