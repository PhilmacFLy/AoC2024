package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type puzzle [][]string

func LoadPuzzle(filename string) puzzle {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var p puzzle
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var line []string
		for _, c := range scanner.Text() {
			line = append(line, string(c))
		}
		p = append(p, line)
	}
	return p
}

func (p puzzle) Print() {
	for _, row := range p {
		for _, c := range row {
			print(c)
		}
		println()
	}
}

func isCountLetter(s string, c int) bool {
	switch c {
	case 0:
		return s == "X"
	case 1:
		return s == "M"
	case 2:
		return s == "A"
	case 3:
		return s == "S"
	}
	return false
}

func (p puzzle) SearchLeft(i, j int) int {
	count := 0
	for j >= 0 && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		j--
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchRight(i, j int) int {
	count := 0
	for j < len(p[i]) && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		j++
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchUp(i, j int) int {
	count := 0
	for i >= 0 && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		i--
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchDown(i, j int) int {
	count := 0
	for i < len(p) && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		i++
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchDiagonalUpLeft(i, j int) int {
	count := 0
	for i >= 0 && j >= 0 && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		i--
		j--
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchDiagonalUpRight(i, j int) int {
	count := 0
	for i >= 0 && j < len(p[i]) && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		i--
		j++
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchDiagonalDownLeft(i, j int) int {
	count := 0
	for i < len(p) && j >= 0 && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		i++
		j--
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) SearchDiagonalDownRight(i, j int) int {
	count := 0
	for i < len(p) && j < len(p[i]) && count < 4 {
		if !isCountLetter(p[i][j], count) {
			return 0
		}
		count++
		i++
		j++
	}
	if count == 4 {
		return 1
	}
	return 0
}

func (p puzzle) Search() int {
	result := 0
	for i, row := range p {
		for j, c := range row {
			if c == "X" {
				result += p.SearchLeft(i, j)
				result += p.SearchRight(i, j)
				result += p.SearchUp(i, j)
				result += p.SearchDown(i, j)
				result += p.SearchDiagonalUpLeft(i, j)
				result += p.SearchDiagonalUpRight(i, j)
				result += p.SearchDiagonalDownLeft(i, j)
				result += p.SearchDiagonalDownRight(i, j)
			}
		}
	}
	return result
}

func getCell(p puzzle, i, j int) string {
	if i < 0 || i >= len(p) || j < 0 || j >= len(p[i]) {
		return ""
	}
	return p[i][j]
}

func upCountIfLetterisCorrect(s string, mcount, scount int) (int, int) {
	if s == "M" {
		mcount += 1
	}
	if s == "S" {
		scount += 1
	}
	return mcount, scount
}

func (p puzzle) FindXMAS(i, j int) int {
	mcount := 0
	scount := 0
	mcount, scount = upCountIfLetterisCorrect(getCell(p, i+1, j+1), mcount, scount)
	mcount, scount = upCountIfLetterisCorrect(getCell(p, i-1, j+1), mcount, scount)
	mcount, scount = upCountIfLetterisCorrect(getCell(p, i+1, j-1), mcount, scount)
	mcount, scount = upCountIfLetterisCorrect(getCell(p, i-1, j-1), mcount, scount)
	diffdiag := true
	if (getCell(p, i+1, j+1) != "") && (getCell(p, i+1, j-1) != "") {
		diffdiag = (getCell(p, i+1, j+1) != getCell(p, i-1, j-1)) && (getCell(p, i-1, j+1) != getCell(p, i+1, j-1))
	}
	if mcount == 2 && scount == 2 && diffdiag {
		return 1
	}
	return 0
}

func (p puzzle) Search2() int {
	result := 0
	for i, row := range p {
		for j, c := range row {
			if c == "A" {
				result += p.FindXMAS(i, j)
			}
		}
	}
	return result
}

func part1(p puzzle) int {
	return p.Search()
}

func part2(p puzzle) int {
	return p.Search2()
}

func main() {
	fmt.Println("Example:")
	p := LoadPuzzle("example.txt")
	fmt.Println("Part1: ", part1(p))
	fmt.Println("Part2: ", part2(p))

	fmt.Println("Input:")
	p = LoadPuzzle("input.txt")
	fmt.Println("Part1: ", part1(p))
	fmt.Println("Part2: ", part2(p))
}
