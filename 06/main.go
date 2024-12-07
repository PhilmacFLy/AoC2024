package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type field struct {
	up       bool
	down     bool
	left     bool
	right    bool
	wall     bool
	tempwall bool
}

func (f field) String() string {
	if f.wall {
		return "#"
	}
	if f.tempwall {
		return "O"
	}
	if f.up {
		return "^"
	}
	if f.down {
		return "v"
	}
	if f.left {
		return "<"
	}
	if f.right {
		return ">"
	}
	return "."
}

func (f *field) setDirection(s string) {
	switch s {
	case "^":
		f.up = true
	case "v":
		f.down = true
	case "<":
		f.left = true
	case ">":
		f.right = true
	}
}

func (f field) isDirection(s string) bool {
	switch s {
	case "^":
		return f.up
	case "v":
		return f.down
	case "<":
		return f.left
	case ">":
		return f.right
	}
	return false
}

func (f field) isEmpty() bool {
	return !f.up && !f.down && !f.left && !f.right
}

type puzzle [][]field

func loadFile(filename string) puzzle {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	p := make(puzzle, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]field, 0)
		for _, c := range line {
			f := field{}
			switch c {
			case '^':
				f.up = true
			case 'v':
				f.down = true
			case '<':
				f.left = true
			case '>':
				f.right = true
			case '#':
				f.wall = true
			}
			row = append(row, f)
		}
		p = append(p, row)
	}
	return p
}

func (p puzzle) print() {
	for _, row := range p {
		for _, cell := range row {
			fmt.Print(cell)
		}
		println()
	}
}

func (p puzzle) findStart() (int, int) {
	for y, row := range p {
		for x, cell := range row {
			if cell.up {
				return x, y
			}
		}
	}
	return -1, -1
}

func copyPuzzle(p puzzle) puzzle {
	newP := make(puzzle, len(p))
	for y, row := range p {
		newRow := make([]field, len(row))
		for x, cell := range row {
			newRow[x] = cell
		}
		newP[y] = newRow
	}
	return newP
}

func (p puzzle) isWall(x, y int) bool {
	return p[y][x].wall || p[y][x].tempwall
}

func (p puzzle) isBorder(x, y int) bool {
	if y > len(p)-1 {
		return true
	}
	return x < 0 || y < 0 || x > len(p[y])-1
}

func turnRight(dx, dy int) (int, int) {
	switch {
	case dx == 0 && dy == 1:
		return -1, 0
	case dx == 0 && dy == -1:
		return 1, 0
	case dx == 1 && dy == 0:
		return 0, 1
	case dx == -1 && dy == 0:
		return 0, -1
	}
	return 0, 0
}

func dxdyToString(dx, dy int) string {
	switch {
	case dx == 0 && dy == 1:
		return "v"
	case dx == 0 && dy == -1:
		return "^"
	case dx == 1 && dy == 0:
		return ">"
	case dx == -1 && dy == 0:
		return "<"
	}
	return "?"
}

func part1(p puzzle) int {
	x, y := p.findStart()
	fmt.Println("Start at:", x, y)
	dx, dy := 0, -1
	steps := 0
	dir := "^"
	for {
		if p.isBorder(x+dx, y+dy) {
			p[y][x].setDirection(dir)
			steps++
			break
		}
		for p.isWall(x+dx, y+dy) {
			dx, dy = turnRight(dx, dy)
			dir = dxdyToString(dx, dy)
		}
		p[y][x].setDirection(dir)
		x += dx
		y += dy
		if p[y][x].isEmpty() {
			steps++
		}
	}
	return steps
}

type counter struct {
	value int
	m     sync.Mutex
}

func (c *counter) inc() {
	c.m.Lock()
	c.value++
	c.m.Unlock()
}

func (c *counter) get() int {
	c.m.Lock()
	defer c.m.Unlock()
	return c.value
}

func (p puzzle) CheckforLoopParallel(x, y int, dx, dy int, c *counter, wg *sync.WaitGroup) {
	defer wg.Done()
	if p.CheckforLoop(x, y, dx, dy) {
		c.inc()
	}
}

func (p puzzle) CheckforLoop(x, y int, dx, dy int) bool {
	visited := make(map[[4]int]bool)

	for {
		state := [4]int{x, y, dx, dy}
		if visited[state] {
			// We have encountered the same position and direction again -> loop
			return true
		}
		visited[state] = true

		// Check if next step is outside the puzzle
		if p.isBorder(x+dx, y+dy) {
			return false
		}

		// If there's a wall ahead, turn right until no wall or border
		for {
			if p.isBorder(x+dx, y+dy) {
				// border encountered while looking for a way forward means no loop, just stop
				return false
			}
			if !p.isWall(x+dx, y+dy) {
				// found a free space ahead
				break
			}
			dx, dy = turnRight(dx, dy)
		}

		// Move forward
		x += dx
		y += dy
	}
}

func part2(p puzzle, parallel bool) int {
	startX, startY := p.findStart()
	// Determine the initial direction from the start symbol
	var dx, dy int
	for _, dir := range []string{"^", "v", "<", ">"} {
		if p[startY][startX].isDirection(dir) {
			switch dir {
			case "^":
				dx, dy = 0, -1
			case "v":
				dx, dy = 0, 1
			case "<":
				dx, dy = -1, 0
			case ">":
				dx, dy = 1, 0
			}
		}
	}

	possibilites := 0
	wg := sync.WaitGroup{}
	c := counter{}

	// Iterate over every cell that is not a wall and not the start position
	for y := 0; y < len(p); y++ {
		for x := 0; x < len(p[y]); x++ {
			// Can't place at start or on an existing wall
			if (x == startX && y == startY) || p[y][x].wall {
				continue
			}

			// Make a copy of the puzzle and place a temp wall here
			newp := copyPuzzle(p)
			newp[y][x].tempwall = true

			// Check if this causes a loop
			if !parallel {
				if newp.CheckforLoop(startX, startY, dx, dy) {
					possibilites++
				}
			} else {
				wg.Add(1)
				go newp.CheckforLoopParallel(startX, startY, dx, dy, &c, &wg)
			}
		}
	}

	if parallel {
		wg.Wait()
		possibilites = c.get()
	}
	return possibilites
}

func main() {
	p := loadFile("input.txt")
	start := time.Now()
	fmt.Println("Part2", part2(p, false))
	fmt.Println("Duration:", time.Since(start))
}
