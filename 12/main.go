package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type plot struct {
	plant  string
	region int
}

type region struct {
	fields    int
	perimeter int
	sides     int
}

func (r region) String() string {
	return fmt.Sprintf("Fields: %d, perimeter: %d, Sides: %d", r.fields, r.perimeter, r.sides)
}

type garden [][]plot

func (g garden) String() string {
	res := ""
	for _, r := range g {
		for _, c := range r {
			res += c.plant
		}
		res += "\n"
	}
	return res
}

func (g garden) Copy() garden {
	var newg garden
	for _, r := range g {
		var newr []plot
		for _, c := range r {
			newr = append(newr, c)
		}
		newg = append(newg, newr)
	}
	return newg
}

func (g garden) WritetoFile(name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range g {
		for _, c := range r {
			f.WriteString(c.plant)
		}
		f.WriteString("\n")
	}
}

func (g garden) PrintRegions() {
	for _, r := range g {
		for _, c := range r {
			r := strconv.Itoa(c.region)
			if len(r) == 1 {
				fmt.Print("0")
			}
			fmt.Print(r + " ")
		}
		fmt.Println()
		fmt.Println()
	}
}

func loadGarden(filename string) garden {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var g garden
	for scanner.Scan() {
		var row []plot
		for _, c := range scanner.Text() {
			row = append(row, plot{plant: string(c)})
		}
		g = append(g, row)
	}
	return g
}

func markRegion(g *garden, x, y int, currentregion int, currentplant string) bool {
	if x < 0 || x >= len((*g)[0]) || y < 0 || y >= len(*g) {
		return false
	}
	if (*g)[y][x].region != 0 || (*g)[y][x].plant != currentplant {
		return false
	}
	(*g)[y][x].region = currentregion
	markRegion(g, x+1, y, currentregion, currentplant)
	markRegion(g, x-1, y, currentregion, currentplant)
	markRegion(g, x, y+1, currentregion, currentplant)
	markRegion(g, x, y-1, currentregion, currentplant)
	return true
}

func countRegionFieldsAndperimeter(g garden, currentregion int) region {
	var r region
	for y, row := range g {
		for x, p := range row {
			if p.region == currentregion {
				r.fields++
				if x+1 >= len(row) || g[y][x+1].region != currentregion {
					r.perimeter++
				}
				if x-1 < 0 || g[y][x-1].region != currentregion {
					r.perimeter++
				}
				if y+1 >= len(g) || g[y+1][x].region != currentregion {
					r.perimeter++
				}
				if y-1 < 0 || g[y-1][x].region != currentregion {
					r.perimeter++
				}
			}
		}
	}
	return r
}

func (g *garden) GetRegion(x, y int) int {
	if x < 0 || x >= len((*g)[0]) || y < 0 || y >= len(*g) {
		return -1
	}
	return (*g)[y][x].region
}

func countCorners(g garden, x, y int) int {
	this := g[y][x].region
	top := this == g.GetRegion(x, y-1)
	right := this == g.GetRegion(x+1, y)
	bottom := this == g.GetRegion(x, y+1)
	left := this == g.GetRegion(x-1, y)
	topright := this == g.GetRegion(x+1, y-1)
	bottomright := this == g.GetRegion(x+1, y+1)
	bottomleft := this == g.GetRegion(x-1, y+1)
	topleft := this == g.GetRegion(x-1, y-1)

	if top && right && bottom && left {
		count := 0
		if !topright {
			count++
		}
		if !bottomright {
			count++
		}
		if !bottomleft {
			count++
		}
		if !topleft {
			count++
		}
		return count
	}

	if top && right && bottom && !left {
		count := 0
		if !topright {
			count++
		}
		if !bottomright {
			count++
		}
		return count
	}

	if top && !right && bottom && left {
		count := 0
		if !bottomleft {
			count++
		}
		if !topleft {
			count++
		}
		return count
	}

	if top && right && !bottom && left {
		count := 0
		if !topright {
			count++
		}
		if !topleft {
			count++
		}
		return count
	}

	if !top && right && bottom && left {
		count := 0
		if !bottomleft {
			count++
		}
		if !bottomright {
			count++
		}
		return count
	}

	if top && !right && !bottom && left {
		count := 1
		if !topleft {
			count++
		}
		return count
	}

	if top && right && !bottom && !left {
		count := 1
		if !topright {
			count++
		}
		return count
	}

	if !top && right && bottom && !left {
		count := 1
		if !bottomright {
			count++
		}
		return count
	}

	if !top && !right && bottom && left {
		count := 1
		if !bottomleft {
			count++
		}
		return count
	}

	if top && !right && !bottom && !left {
		return 2
	}

	if !top && right && !bottom && !left {
		return 2
	}

	if !top && !right && bottom && !left {
		return 2
	}

	if !top && !right && !bottom && left {
		return 2
	}

	if !top && !right && !bottom && !left {
		return 4
	}

	return 0

}

func part1(g garden) int {
	currentregion := 1
	for y, r := range g {
		for x := range r {
			if markRegion(&g, x, y, currentregion, g[y][x].plant) {
				currentregion++
			}
		}
	}
	var regions []region
	for i := 1; i < currentregion; i++ {
		r := countRegionFieldsAndperimeter(g, i)
		regions = append(regions, r)
	}
	result := 0
	for _, r := range regions {
		result += r.fields * r.perimeter
	}
	return result
}

func part2(g garden) int {
	currentregion := 1
	for y, r := range g {
		for x := range r {
			if markRegion(&g, x, y, currentregion, g[y][x].plant) {
				currentregion++
			}
		}
	}
	var regions []region
	for i := 1; i < currentregion; i++ {
		r := countRegionFieldsAndperimeter(g, i)
		regions = append(regions, r)
	}
	for y, r := range g {
		for x := range r {
			corners := countCorners(g, x, y)
			regions[g[y][x].region-1].sides += corners
		}
	}

	result := 0
	for _, r := range regions {
		fmt.Println(r)
		result += r.fields * r.sides
	}

	return result
}

func main() {
	fmt.Println("Part1 Example:", part1(loadGarden("example.txt")))
	fmt.Println("Part1 Input:", part1(loadGarden("input.txt")))
	fmt.Println("Part2 Example:", part2(loadGarden("example.txt")))
	fmt.Println("Part2 Input:", part2(loadGarden("input.txt")))
}
