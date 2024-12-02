package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func importInput(filename string) ([]int, []int) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	left := make([]int, 0)
	right := make([]int, 0)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		if len(split) != 4 {
			log.Fatal("Length of input is not 2 but ", len(split))
			log.Fatal("Invalid input")
		}
		x, err := strconv.Atoi(split[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(split[3])
		if err != nil {
			log.Fatal(err)
		}
		left = append(left, x)
		right = append(right, y)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return left, right
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func part1(left, right []int) {
	distance := 0
	for i := 0; i < len(left); i++ {
		distance += abs(right[i] - left[i])
	}
	fmt.Println("Distance is", distance)
}

func part2(left, right []int) {
	occurence := make(map[int]int)
	for _, x := range right {
		occurence[x]++
	}
	similar := 0
	for _, x := range left {
		if occurence[x] > 0 {
			similar += x * occurence[x]
		}
	}
	fmt.Println("Similar is", similar)
}

func main() {
	left, right := importInput("input.txt")
	sort.Ints(left)
	sort.Ints(right)
	part1(left, right)
	part2(left, right)
}
