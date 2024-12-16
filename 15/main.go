package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type field [][]string

func (f field) String() string {
	var s string
	for _, row := range f {
		s += strings.Join(row, "") + "\n"
	}
	return s
}

type dir int

const (
	up dir = iota
	down
	left
	right
)

type orders []dir

func (o orders) String() string {
	var s string
	for _, d := range o {
		switch d {
		case up:
			s += "^"
		case down:
			s += "v"
		case left:
			s += "<"
		case right:
			s += ">"
		}
	}
	return s
}

func LoadFromFiled(filename string) (field, orders) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var fi field
	var o orders
	scanner := bufio.NewScanner(f)
	mode := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mode++
			continue
		}
		if mode == 0 {
			fi = append(fi, strings.Split(line, ""))
		} else {
			for _, c := range line {
				switch c {
				case '^':
					o = append(o, up)
				case 'v':
					o = append(o, down)
				case '<':
					o = append(o, left)
				case '>':
					o = append(o, right)
				}
			}
		}
	}
	return fi, o
}

func GetDxDy(d dir) (int, int) {
	switch d {
	case up:
		return 0, -1
	case down:
		return 0, 1
	case left:
		return -1, 0
	case right:
		return 1, 0
	}
	return 0, 0
}

func (f field) CanMove(x, y int, direction dir) (int, int, bool) {
	dx, dy := GetDxDy(direction)
	for {
		x += dx
		y += dy
		if x < 0 || y < 0 || x >= len(f[0]) || y >= len(f) || f[y][x] == "#" {
			return x, y, false
		}
		if f[y][x] == "." {
			return x, y, true
		}
	}
}

func (f *field) Push(fromx, fromy, tox, toy int, direction dir) {
	dx, dy := GetDxDy(direction)
	dx = dx * -1
	dy = dy * -1
	x, y := tox, toy
	if fromx == tox && fromy == toy {
		return
	}
	for {
		(*f)[y][x] = (*f)[y+dy][x+dx]
		x += dx
		y += dy
		if x == fromx && y == fromy {
			break
		}
	}
}
func (f *field) Move(o orders, startx, starty int) {
	x, y := startx, starty
	for _, d := range o {
		tox, toy, ok := f.CanMove(x, y, d)
		if ok {
			dx, dy := GetDxDy(d)
			f.Push(x+dx, y+dy, tox, toy, d)
			(*f)[y][x] = "."
			x = x + dx
			y = y + dy
			(*f)[y][x] = "@"
		}
	}
}

func (f field) FindRobot() (int, int) {
	for y, row := range f {
		for x, c := range row {
			if c == "@" {
				return x, y
			}
		}
	}
	return -1, -1
}

func (f field) SumUpBoxCoordinates() int {
	result := 0
	for y, row := range f {
		for x, c := range row {
			if c == "O" {
				result += x + y*100
			}
		}
	}
	return result
}

func main() {
	f, o := LoadFromFiled("smallexample.txt")
	x, y := f.FindRobot()
	f.Move(o, x, y)
	fmt.Println("Part1 Small Example: ", f.SumUpBoxCoordinates())
	f, o = LoadFromFiled("example.txt")
	x, y = f.FindRobot()
	f.Move(o, x, y)
	fmt.Println("Part1 Example: ", f.SumUpBoxCoordinates())
	f, o = LoadFromFiled("input.txt")
	x, y = f.FindRobot()
	f.Move(o, x, y)
	fmt.Println("Part1: ", f.SumUpBoxCoordinates())
}
