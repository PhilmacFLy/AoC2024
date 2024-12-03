package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
)

type muls struct {
	x, y    int
	enabled bool
}

var regexpart1 = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
var regexpart2 = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

func loadMuls(filename string) []muls {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	input := string(b)
	m := make([]muls, 0)
	matches := regexpart1.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) == 3 {
			x, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			m = append(m, muls{x, y, true})
		}
	}
	return m
}

func loadEnabledMuls(filename string) []muls {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	input := string(b)
	m := make([]muls, 0)
	matches := regexpart2.FindAllStringSubmatch(input, -1)
	enabled := true
	for _, match := range matches {
		if len(match) > 1 && match[1] != "" && match[2] != "" {
			x, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			m = append(m, muls{x, y, enabled})
		} else {
			if match[0] == "do()" {
				enabled = true
			} else if match[0] == "don't()" {
				enabled = false
			}
		}
	}
	return m
}

func part1(m []muls) {
	result := 0
	for _, mul := range m {
		result += mul.x * mul.y
	}
	log.Println("Part 1:", result)
}

func part2(m []muls) {
	result := 0
	for _, mul := range m {
		if mul.enabled {
			result += mul.x * mul.y
		}
	}
	log.Println("Part 2:", result)
}

func main() {
	log.Println("========== Part 1 ==========")
	log.Println("Example:")
	m := loadMuls("example.txt")
	part1(m)
	log.Println("Input:")
	m = loadMuls("input.txt")
	part1(m)
	log.Println("========== Part 2 ==========")
	log.Println("Example:")
	m = loadEnabledMuls("example2.txt")
	part2(m)
	log.Println("Input:")
	m = loadEnabledMuls("input.txt")
	part2(m)
}
