package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type robot struct {
	Id   int
	PosX int
	PosY int
	VelX int
	VelY int
}

const (
	examplewidth  = 11
	exampleheight = 7
	realwidth     = 101
	realheight    = 103
)

func (r robot) String() string {
	return fmt.Sprintf("Robot: Pos %v,%v, Vel %v %v", r.PosX, r.PosY, r.VelX, r.VelY)
}

func (r *robot) Move(example bool) {
	r.PosX += r.VelX
	r.PosY += r.VelY
	if example {
		r.PosX = (r.PosX + examplewidth) % examplewidth
		r.PosY = (r.PosY + exampleheight) % exampleheight
	} else {
		r.PosX = (r.PosX + realwidth) % realwidth
		r.PosY = (r.PosY + realheight) % realheight
	}
}

func LoadRobotsFromFile(filename string) []robot {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var robots []robot
	scanner := bufio.NewScanner(f)
	count := 1
	for scanner.Scan() {
		line := scanner.Text()
		var r robot
		r.Id = count
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &r.PosX, &r.PosY, &r.VelX, &r.VelY)
		robots = append(robots, r)
		count++
	}
	return robots
}

func PrintField(r []robot, example bool) {
	var width, height int
	if example {
		width = examplewidth
		height = exampleheight
	} else {
		width = realwidth
		height = realheight
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			count := 0
			for _, robot := range r {
				if robot.PosX == x && robot.PosY == y {
					count++
				}
			}
			if count > 0 {
				fmt.Print(count)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func countRobotsinQuadrants(r []robot, example bool) int {
	var width, height int
	if example {
		width = examplewidth
		height = exampleheight
	} else {
		width = realwidth
		height = realheight
	}
	quadrants := make([]int, 4)
	PrintField(r, example)
	for _, robot := range r {
		fmt.Printf("Robot Pos: %v %v Width Height %v %v \n", robot.PosX, robot.PosY, width/2, height/2)
		if robot.PosX == width/2 || robot.PosY == height/2 {
			continue
		}
		if robot.PosX < width/2 && robot.PosY < height/2 {
			quadrants[0]++
		} else if robot.PosX >= width/2 && robot.PosY < height/2 {
			quadrants[1]++
		} else if robot.PosX < width/2 && robot.PosY >= height/2 {
			quadrants[2]++
		} else if robot.PosX >= width/2 && robot.PosY >= height/2 {
			quadrants[3]++
		}
	}
	fmt.Println(quadrants)
	result := 1
	for _, q := range quadrants {
		if q > 0 {
			result *= q
		}
	}
	return result
}

func findTenRobtsInLine(r []robot, example bool) bool {
	var width, height int
	if example {
		width = examplewidth
		height = exampleheight
	} else {
		width = realwidth
		height = realheight
	}
	var field [][]string
	for y := 0; y < height; y++ {
		var row []string
		for x := 0; x < width; x++ {
			count := 0
			for _, robot := range r {
				if robot.PosX == x && robot.PosY == y {
					count++
				}
			}
			if count > 0 {
				row = append(row, "#")
			} else {
				row = append(row, ".")
			}
		}
		field = append(field, row)
	}
	for y := 0; y < height-1; y++ {
		for x := 0; x < width-1; x++ {
			if field[y][x] == "#" && (x+9) < width {
				if field[y][x+1] == "#" && field[y][x+2] == "#" && field[y][x+3] == "#" && field[y][x+4] == "#" && field[y][x+5] == "#" && field[y][x+6] == "#" && field[y][x+7] == "#" && field[y][x+8] == "#" && field[y][x+9] == "#" {
					return true
				}
			}
		}
	}
	return false
}

func part1(r []robot, example bool) int {
	for i := 0; i < 100; i++ {
		for j := range r {
			r[j].Move(example)
		}
	}
	return countRobotsinQuadrants(r, example)
}

func part2(r []robot, example bool) {
	var fields string
	for i := 0; i < 100000; i++ {
		for j := range r {
			r[j].Move(example)
		}
		if findTenRobtsInLine(r, example) {
			fmt.Println("Found at", i)
			PrintField(r, example)
		}

	}
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(fields)
}

func main() {
	fmt.Println("Part 1 Example:", part1(LoadRobotsFromFile("example.txt"), true))
	fmt.Println("Part 1:", part1(LoadRobotsFromFile("input.txt"), false))
	part2(LoadRobotsFromFile("input.txt"), false)
}
