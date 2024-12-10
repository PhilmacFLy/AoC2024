package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type pos struct {
	height  int
	reached bool
}

type puzzle [][]pos

func (p puzzle) String() string {
	res := ""
	for _, r := range p {
		for _, c := range r {
			if c.height == -1 {
				res += "."
			} else {
				res += strconv.Itoa(c.height)
			}
		}
		res += "\n"
	}
	return res
}

func (p puzzle) PrintReached() {
	for _, r := range p {
		for _, c := range r {
			if c.reached {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (p *puzzle) ResetReached() {
	for y, r := range *p {
		for x := range r {
			(*p)[y][x].reached = false
		}
	}
}

func (p *puzzle) GetValue(x, y int) int {
	if x < 0 || y < 0 || x >= len(*p) || y >= len((*p)[0]) {
		return -1
	}
	return (*p)[y][x].height
}

func (p *puzzle) LeftValue(x, y int) int {
	x -= 1
	return p.GetValue(x, y)
}

func (p *puzzle) RightValue(x, y int) int {
	x += 1
	return p.GetValue(x, y)
}

func (p *puzzle) UpValue(x, y int) int {
	y -= 1
	return p.GetValue(x, y)
}

func (p *puzzle) DownValue(x, y int) int {
	y += 1
	return p.GetValue(x, y)
}

func (p *puzzle) SetReached(x, y int) {
	(*p)[y][x].reached = true
}

func LoadFromFile(filename string) puzzle {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var p puzzle
	for scanner.Scan() {
		var row []pos
		for _, c := range scanner.Text() {
			s := 0
			if c == '.' {
				s = -1
			} else {
				s, _ = strconv.Atoi(string(c))
			}
			row = append(row, pos{height: s, reached: false})
		}
		p = append(p, row)
	}
	return p
}

type node struct {
	x, y int
}

func (p puzzle) WideSearch(x, y int) int {
	queue := make([]node, 0)
	queue = append(queue, node{x, y})
	p.SetReached(x, y)
	res := 0
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		thisval := p.GetValue(n.x, n.y)
		if thisval == -1 {
			continue
		}

		if (thisval == 9) && !p[n.y][n.x].reached {
			res++
			p.SetReached(n.x, n.y)
			continue
		}
		if p.LeftValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x - 1, n.y})
			p.SetReached(n.x, n.y)
		}
		if p.RightValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x + 1, n.y})
			p.SetReached(n.x, n.y)
		}
		if p.UpValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x, n.y - 1})
			p.SetReached(n.x, n.y)
		}
		if p.DownValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x, n.y + 1})
			p.SetReached(n.x, n.y)
		}
	}
	return res
}

func (p puzzle) WideSearch2(x, y int) int {
	queue := make([]node, 0)
	queue = append(queue, node{x, y})
	res := 0
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		thisval := p.GetValue(n.x, n.y)
		if thisval == -1 {
			continue
		}

		if (thisval == 9) && !p[n.y][n.x].reached {
			res++
			continue
		}
		if p.LeftValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x - 1, n.y})
		}
		if p.RightValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x + 1, n.y})
		}
		if p.UpValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x, n.y - 1})
		}
		if p.DownValue(n.x, n.y) == thisval+1 {
			queue = append(queue, node{n.x, n.y + 1})
		}
	}
	return res
}

func part1(p puzzle) int {
	fmt.Println(p)
	res := 0
	for y, r := range p {
		for x, c := range r {
			if c.height == 0 {
				p.ResetReached()
				locares := p.WideSearch(x, y)
				res += locares
			}
		}
	}

	return res
}

func part2(p puzzle) int {
	fmt.Println(p)
	res := 0
	for y, r := range p {
		for x, c := range r {
			if c.height == 0 {
				locares := p.WideSearch2(x, y)
				res += locares
			}
		}
	}

	return res
}

func main() {
	log.Println("Part1 Example:", part1(LoadFromFile("example.txt")))
	log.Println("Part1 Input:", part1(LoadFromFile("input.txt")))
	log.Println("Part2 Example:", part2(LoadFromFile("example.txt")))
	log.Println("Part2 Input:", part2(LoadFromFile("input.txt")))
}
