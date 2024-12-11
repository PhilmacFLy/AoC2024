package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type stones []int

type stone struct {
	valueleft  int
	valueright int
}

type cache map[int]stone

var c cache

func LoadFromFile(filename string) stones {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var p stones
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		for _, c := range split {
			i, _ := strconv.Atoi(c)
			p = append(p, i)
		}
	}
	return p
}

func (s *stones) String() string {
	var str strings.Builder
	for _, i := range *s {
		str.WriteString(strconv.Itoa(i) + " ")
	}
	return str.String()
}

func (s *stones) InsertValue(i, v int) {
	*s = append((*s)[:i], append([]int{v}, (*s)[i:]...)...)
}

func (s *stones) RemoveValue(i int) {
	*s = append((*s)[:i], (*s)[i+1:]...)
}

func (s *stones) Blink() {
	i := 0
	for i < len(*s) {
		val := (*s)[i]
		if cache, ok := c[val]; ok {
			(*s)[i] = cache.valueleft
			i += 1
			if cache.valueright != -1 {
				s.InsertValue(i, cache.valueright)
				i += 1
			}
		} else {
			if val == 0 {
				(*s)[i] = 1
				i += 1
			} else {
				is := strconv.Itoa(val)
				if len(is)%2 == 0 {
					left := is[:len(is)/2]
					right := is[len(is)/2:]
					lefti, _ := strconv.Atoi(left)
					righti, _ := strconv.Atoi(right)
					(*s)[i] = lefti
					s.InsertValue(i+1, righti)
					c[val] = stone{lefti, righti}
					i += 2
				} else {
					(*s)[i] = (*s)[i] * 2024
					c[val] = stone{(*s)[i], -1}
					i += 1
				}
			}
		}
	}
}

func blinking(times int, s stones) int {
	for i := 0; i < times; i += 1 {
		log.Println("Blinking", i)
		s.Blink()
	}
	return len(s)
}

func blinking2(times int, s stones) int {
	stonecount := make(map[int]int)
	for _, st := range s {
		stonecount[st] += 1
	}

	for i := 0; i < times; i += 1 {
		newStoneCount := make(map[int]int)
		for stone, count := range stonecount {
			var createdStones []int
			if stone == 0 {
				createdStones = append(createdStones, 1)
			} else if digits := strconv.Itoa(stone); len(digits)%2 == 0 {
				a, _ := strconv.Atoi(digits[:len(digits)/2])
				b, _ := strconv.Atoi(digits[len(digits)/2:])
				createdStones = append(createdStones, a)
				createdStones = append(createdStones, b)
			} else {
				createdStones = append(createdStones, stone*2024)
			}
			for _, newStone := range createdStones {
				newStoneCount[newStone] = newStoneCount[newStone] + count
			}
		}
		stonecount = newStoneCount
	}
	result := 0
	for _, count := range stonecount {
		result += count
	}
	return result
}

func main() {
	c = make(cache)
	log.Println("Part 1 Example:", blinking2(25, LoadFromFile("example.txt")))
	log.Println("Part 1:", blinking2(25, LoadFromFile("input.txt")))
	log.Println("Part 2 Example:", blinking2(75, LoadFromFile("example.txt")))
	log.Println("Part 2:", blinking2(5000, LoadFromFile("input.txt")))
}
