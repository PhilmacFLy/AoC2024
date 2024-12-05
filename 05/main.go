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

type rules struct {
	isSmallerThan []int
	isGreaterThan []int
}

func (r rules) checkIfSmallerThan(x int) bool {
	for _, val := range r.isSmallerThan {
		if val == x {
			return true
		}
	}
	return false
}

func (r rules) checkIfGreaterThan(x int) bool {
	for _, val := range r.isGreaterThan {
		if val == x {
			return true
		}
	}
	return false
}

type ruleset map[int]rules
type updates [][]int

func loadFile(filename string) (ruleset, updates) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	loadingrulsets := true
	ruleset := make(ruleset)
	updates := make(updates, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			loadingrulsets = false
			continue
		}
		if loadingrulsets {
			split := strings.Split(line, "|")
			if len(split) != 2 {
				log.Fatal("invalid rule")
			}
			x, err := strconv.Atoi(split[0])
			if err != nil {
				log.Fatalf("invalid number in rules: %s", split[0])
			}
			y, err := strconv.Atoi(split[1])
			if err != nil {
				log.Fatalf("invalid number in rules: %s", split[1])
			}
			if _, ok := ruleset[x]; !ok {
				ruleset[x] = rules{}
			}
			rule := ruleset[x]
			rule.isSmallerThan = append(rule.isSmallerThan, y)
			ruleset[x] = rule
			if _, ok := ruleset[y]; !ok {
				ruleset[y] = rules{}
			}
			rule = ruleset[y]
			rule.isGreaterThan = append(rule.isGreaterThan, x)
			ruleset[y] = rule
		} else {
			var update []int
			for _, val := range strings.Split(line, ",") {
				x, err := strconv.Atoi(val)
				if err != nil {
					log.Fatalf("invalid number in updates: %s", val)
				}
				update = append(update, x)
			}
			updates = append(updates, update)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ruleset, updates
}

func checkForward(r rules, forward []int) bool {
	result := true
	for _, f := range forward {
		if !r.checkIfSmallerThan(f) {
			result = false
		}
	}
	return result
}

func checkBackward(r rules, backward []int) bool {
	result := true
	for i := len(backward) - 1; i >= 0; i-- {
		if !r.checkIfGreaterThan(backward[i]) {
			result = false
		}
	}
	return result
}

func part1(r ruleset, u updates) int {
	var valid [][]int
	for _, update := range u {
		isvalid := true
		for i, val := range update {
			var forward []int
			var backward []int
			if i > 0 {
				backward = update[:i]
			}
			if i < len(update)-1 {
				forward = update[i+1:]
			}
			if !checkForward(r[val], forward) || !checkBackward(r[val], backward) {
				isvalid = false
			}
		}
		if isvalid {
			valid = append(valid, update)
		}
	}
	result := 0
	for _, v := range valid {
		if len(v)%2 == 0 {
			fmt.Println("Even number of elements")
		}
		pos := len(v) / 2
		result += v[pos]
	}
	return result
}

func sortByRules(ruleset ruleset, nums []int) []int {
	sort.Slice(nums, func(i, j int) bool {
		x, y := nums[i], nums[j]

		// Check if there's a direct rule for x < y
		for _, smaller := range ruleset[x].isSmallerThan {
			if smaller == y {
				return true
			}
		}
		// Check if there's a direct rule for y < x
		for _, greater := range ruleset[y].isGreaterThan {
			if greater == x {
				return false
			}
		}

		// Check the other way around
		for _, greater := range ruleset[x].isGreaterThan {
			if greater == y {
				return false
			}
		}
		for _, smaller := range ruleset[y].isSmallerThan {
			if smaller == x {
				return true
			}
		}

		// Fallback to natural order if no rule is specified
		panic("no rule specified")
	})
	return nums
}

func part2(r ruleset, u updates) int {
	var invalid [][]int
	for _, update := range u {
		isvalid := true
		for i, val := range update {
			var forward []int
			var backward []int
			if i > 0 {
				backward = update[:i]
			}
			if i < len(update)-1 {
				forward = update[i+1:]
			}
			if !checkForward(r[val], forward) || !checkBackward(r[val], backward) {
				isvalid = false
			}
		}
		if !isvalid {
			invalid = append(invalid, update)
		}
	}
	// Resort the invalid list
	result := 0
	for i := 0; i < len(invalid); i++ {
		invalid[i] = sortByRules(r, invalid[i])
	}
	for _, v := range invalid {
		if len(v)%2 == 0 {
			fmt.Println("Even number of elements")
		}
		pos := len(v) / 2
		result += v[pos]
	}
	return result

}

func main() {
	fmt.Println("Example:")
	r, u := loadFile("example.txt")
	fmt.Println("Part 1: ", part1(r, u))
	fmt.Println("Part 2: ", part2(r, u))

	fmt.Println("Input:")
	r, u = loadFile("input.txt")
	fmt.Println("Part 1: ", part1(r, u))
	fmt.Println("Part 2: ", part2(r, u))
}
