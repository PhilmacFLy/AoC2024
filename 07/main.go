package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type calibration struct {
	result int
	parts  []int
}

func loadCalibrations(filename string) []calibration {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var calibrations []calibration
	for scanner.Scan() {
		line := scanner.Text()
		var c calibration
		split := strings.Split(line, ":")
		c.result, err = strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}
		parts := strings.Split(strings.TrimSpace(split[1]), " ")
		for _, part := range parts {
			partInt, err := strconv.Atoi(part)
			if err != nil {
				panic(err)
			}
			c.parts = append(c.parts, partInt)
		}
		calibrations = append(calibrations, c)
	}
	return calibrations
}

func (c calibration) String() string {
	return strconv.Itoa(c.result) + ":" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(c.parts)), " "), "[]")
}

func calculateOptions(nums []int, index int, current int, result int) bool {
	// If we've processed all numbers, add the result to the list
	if index == len(nums) {
		return current == result
	}

	// Add the current number
	if calculateOptions(nums, index+1, current+nums[index], result) {
		return true
	}

	// Multiply by the current number
	return calculateOptions(nums, index+1, current*nums[index], result)
}

func (c calibration) isValid() bool {
	return calculateOptions(c.parts, 1, c.parts[0], c.result)
}

func calculateOptions2(nums []int, index int, current int, result int) bool {
	// If we've processed all numbers, add the result to the list
	if index == len(nums) {
		return current == result
	}

	// Add the current number
	if calculateOptions2(nums, index+1, current+nums[index], result) {
		return true
	}

	// Multiply by the current number
	if calculateOptions2(nums, index+1, current*nums[index], result) {
		return true
	}

	// Concatenate the current number
	concatNum := strconv.Itoa(current) + strconv.Itoa(nums[index])
	num, _ := strconv.Atoi(concatNum)
	return calculateOptions2(nums, index+1, num, result)
}

func (c calibration) isValid2() bool {
	return calculateOptions2(c.parts, 1, c.parts[0], c.result)
}

type counter struct {
	value int
	m     sync.Mutex
}

func (c *counter) inc(i int) {
	c.m.Lock()
	c.value += i
	c.m.Unlock()
}

func (c *counter) get() int {
	c.m.Lock()
	defer c.m.Unlock()
	return c.value
}

func (c calibration) isValidParallel(wg *sync.WaitGroup, counter *counter) {
	defer wg.Done()
	if c.isValid() {
		counter.inc(c.result)
	}
}

func (c calibration) isValid2Parallel(wg *sync.WaitGroup, counter *counter) {
	defer wg.Done()
	result := c.isValid2()
	if result {
		counter.inc(c.result)
	}
	fmt.Println("Result:", c.result, "Valid:", result)
}

func part2(calibrations []calibration) int {
	result := 0
	wg := sync.WaitGroup{}
	counter := counter{}
	for _, c := range calibrations {
		wg.Add(1)
		go c.isValid2Parallel(&wg, &counter)
	}
	wg.Wait()
	result = counter.get()
	return result
}

func part1(calibrations []calibration) int {
	result := 0
	wg := sync.WaitGroup{}
	counter := counter{}
	for _, c := range calibrations {
		wg.Add(1)
		go c.isValidParallel(&wg, &counter)
	}
	wg.Wait()
	result = counter.get()
	return result
}

func main() {
	calibrations := loadCalibrations("example.txt")
	fmt.Println("Part 1 Example:", part1(calibrations))
	calibrations = loadCalibrations("input.txt")
	fmt.Println("Part 1:", part1(calibrations))
	calibrations = loadCalibrations("example.txt")
	fmt.Println("Part 2 Example:", part2(calibrations))
	calibrations = loadCalibrations("input.txt")
	fmt.Println("Part 2:", part2(calibrations))
}
