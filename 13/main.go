package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type equ struct {
	a []int
	b []int
	x int
	y int
}

func (e equ) String() string {
	result := fmt.Sprintf("Button A: %v %v\n", e.a[0], e.a[1])
	result += fmt.Sprintf("Button B: %v %v\n", e.b[0], e.b[1])
	result += fmt.Sprintf("Prize: X=%v, Y=%v", e.x, e.y)
	return result
}

func getAB(input string) []int {
	sab := strings.Split(input, ",")
	var ab []int
	for _, v := range sab {
		v = strings.TrimSpace(v)
		v = strings.ReplaceAll(v, "X", "")
		v = strings.ReplaceAll(v, "Y", "")
		temp, _ := strconv.Atoi(v)
		ab = append(ab, temp)
	}
	return ab
}

func getXY(input string) (int, int) {
	sxy := strings.Split(input, ",")
	var xy [2]int
	for i, v := range sxy {
		v = strings.TrimSpace(v)
		v = strings.ReplaceAll(v, "X", "")
		v = strings.ReplaceAll(v, "Y", "")
		v = strings.ReplaceAll(v, "=", "")
		xy[i], _ = strconv.Atoi(v)
	}
	return xy[0], xy[1]
}

func roundToPrecision(value float64, precision int) float64 {
	scale := math.Pow(10, float64(precision))
	return math.Round(value*scale) / scale
}

func HasDecimal(value float64) bool {
	return value != float64(int(value))
}

func (e equ) solve() (float64, float64) {
	variables := []float64{
		float64(e.a[0]), float64(e.b[0]),
		float64(e.a[1]), float64(e.b[1]),
	}
	variableVec := mat.NewDense(2, 2, variables)

	// Define the result vector b
	result := []float64{float64(e.x), float64(e.y)}
	resVec := mat.NewVecDense(2, result)

	// Solve the system of equations Ax = b
	var x mat.VecDense
	err := x.SolveVec(variableVec, resVec)
	if err != nil {
		log.Fatalf("Failed to solve the system: %v", err)
	}

	return roundToPrecision(x.At(0, 0), 2), roundToPrecision(x.At(1, 0), 2)
}

func LoadEquations(filename string) []equ {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var equations []equ
	scanner := bufio.NewScanner(f)
	counter := 0
	var equation equ
	for scanner.Scan() {
		mod := counter % 4
		if mod == 3 {
			counter++
			equations = append(equations, equation)
			equation = equ{}
			continue
		}

		line := strings.Split(scanner.Text(), ":")
		work := strings.TrimSpace(line[1])
		switch mod {
		case 0:
			equation.a = append(equation.a, getAB(work)...)
		case 1:
			equation.b = append(equation.b, getAB(work)...)
		case 2:
			equation.x, equation.y = getXY(work)
		}
		counter++
	}
	equations = append(equations, equation)
	return equations
}

func part1(equations []equ) float64 {
	result := 0.0
	for _, e := range equations {
		a, b := e.solve()
		//fmt.Printf("A: %v, B: %v\n", a, b)
		if !HasDecimal(a) && !HasDecimal(b) && a > 0 && b > 0 {
			result += a*3 + b
		}
	}
	return result
}

func part2(equations []equ) float64 {
	for i := range equations {
		equations[i].x += 10000000000000
		equations[i].y += 10000000000000
	}
	return part1(equations)
}

func main() {
	fmt.Println("Part 1 Example:", part1(LoadEquations("example.txt")))
	fmt.Println("Part 1:", part1(LoadEquations("input.txt")))
	fmt.Println("Part 2 Example:", part2(LoadEquations("example.txt")))
	fmt.Println("Part 2:", part2(LoadEquations("input.txt")))
}
