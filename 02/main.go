package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type report []int

func importData(filename string) []report {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	data := make([]report, 0)
	for scanner.Scan() {
		r := make(report, 0)
		split := strings.Split(scanner.Text(), " ")
		for _, s := range split {
			x, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			r = append(r, x)
		}
		data = append(data, r)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isIncreasing(x, y int) bool {
	return x <= y
}

func isAtMostThree(x, y int) bool {
	diff := abs(y - x)
	return (diff > 0 && diff <= 3)
}

func (r report) remove(i int) report {
	newreport := make(report, len(r)-1)
	copy(newreport, r[:i])
	copy(newreport[i:], r[i+1:])
	return newreport
}

func (r report) isSafe() bool {
	if len(r) < 2 {
		return true
	}
	for i := 0; i < len(r)-1; i++ {
		if !isAtMostThree(r[i], r[i+1]) {
			return false
		}
		if len(r) > i+2 {
			if isIncreasing(r[i], r[i+1]) != isIncreasing(r[i+1], r[i+2]) {
				return false
			}
		}
	}
	return true
}

func (r report) isSafeDamper() bool {
	if len(r) < 2 {
		return true
	}
	for i := 0; i < len(r)-1; i++ {
		if !isAtMostThree(r[i], r[i+1]) {
			newreport := r.remove(i)
			if newreport.isSafe() {
				return true
			}
			newreport = r.remove(i + 1)
			return newreport.isSafe()
		}
		if len(r) > i+2 {
			if isIncreasing(r[i], r[i+1]) != isIncreasing(r[i+1], r[i+2]) {
				newreport := r.remove(i)
				if newreport.isSafe() {
					return true
				}
				newreport = r.remove(i + 1)
				if newreport.isSafe() {
					return true
				}
				newreport = r.remove(i + 2)
				return newreport.isSafe()
			}
		}
	}
	return true
}

func part1(reports []report) {
	count := 0
	for _, r := range reports {
		if r.isSafe() {
			count++
		}
	}
	log.Println("Safe reports:", count)
}

func part2(reports []report) {
	count := 0
	for i, r := range reports {
		if r.isSafeDamper() {
			count++
			fmt.Println("Report", i+1, "is safe")
		} else {
			fmt.Println("Report", i+1, "is not safe")
		}
	}
	log.Println("Safe reports with damper:", count)
}

func main() {
	fmt.Println("Example:")
	fmt.Println("===================================")
	reports := importData("example.txt")
	part1(reports)
	part2(reports)

	fmt.Println()
	fmt.Println("Input:")
	fmt.Println("===================================")
	reports = importData("input.txt")
	part1(reports)
	part2(reports)
}
