package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func (p puzzle) CheckforLoop(x, y int, dx, dy int, dir string) bool {
	for {
		p[y][x].setDirection(dir)
		if p.isBorder(x+dx, y+dy) {
			return false
		}
		for p.isWall(x+dx, y+dy) {
			dx, dy = turnRight(dx, dy)
			dir = dxdyToString(dx, dy)
		}
		if !p.isBorder(x+dx, y+dy) {
			if p[y+dy][x+dx].isDirection(dir) {
				return true
			}
			x += dx
			y += dy
		}
	}

}

func part2(p puzzle) int {
	x, y := p.findStart()
	fmt.Println("Start at:", x, y)
	dx, dy := 0, -1
	possibilites := 0
	dir := "^"
	for {
		p[y][x].setDirection(dir)
		if p.isBorder(x+dx, y+dy) {
			p[y][x].setDirection(dir)
			break
		}
		for p.isWall(x+dx, y+dy) {
			dx, dy = turnRight(dx, dy)
			dir = dxdyToString(dx, dy)
		}
		p[y][x].setDirection(dir)
		if !p.isBorder(x+dx, y+dy) && !p.isWall(x+dx, y+dy) {
			newp := copyPuzzle(p)
			newp[x+dx][y+dy].tempwall = true
			tempdx, tempdy := turnRight(dx, dy)
			tempdir := dxdyToString(tempdx, tempdy)
			if newp.CheckforLoop(x, y, tempdx, tempdy, tempdir) {
				possibilites++
			}
		}

		x += dx
		y += dy
	}
	return possibilites
}

func main() {
	fmt.Println("Example:")
	p := loadFile("example.txt")
	fmt.Println("Part1:", part1(p))
	p = loadFile("example.txt")
	fmt.Println("Part2:", part2(p))

	fmt.Println("Input:")
	p = loadFile("input.txt")
	fmt.Println("Part1:", part1(p))
	p = loadFile("input.txt")
	fmt.Println("Part2:", part2(p))
}
