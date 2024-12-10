package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type content struct {
	id     int
	isfile bool
	size   int
}

func (c content) String() string {
	if c.isfile {
		return fmt.Sprintf("File %d Size %d", c.id, c.size)
	}
	return fmt.Sprintf("Free Space %d", c.size)
}

type diskmap []content

func (d diskmap) String() string {
	res := ""
	for _, c := range d {
		p := "."
		if c.isfile {
			p = strconv.Itoa(c.id)
		}
		for i := 0; i < c.size; i++ {
			res += p
		}
	}
	return res
}

func (d *diskmap) Remove(i int) {
	*d = append((*d)[:i], (*d)[i+1:]...)
}

func (d *diskmap) Insert(i int, c content) {
	*d = append(*d, content{})
	copy((*d)[i+1:], (*d)[i:])
	(*d)[i] = c
}

func (d *diskmap) defrag() {
	i := len(*d) - 1
	for i >= 0 {
		if (*d)[i].isfile {
			j := 0
			for j < i {
				if !(*d)[j].isfile {
					if (*d)[j].size == (*d)[i].size {
						(*d)[j].isfile = true
						(*d)[j].id = (*d)[i].id
						(*d)[i].isfile = false
						j++
						break
					} else if (*d)[j].size > (*d)[i].size {
						(*d)[j].size -= (*d)[i].size
						(*d).Insert(j, content{id: (*d)[i].id, isfile: true, size: (*d)[i].size})
						(*d)[i+1].isfile = false
						break
					} else if (*d)[j].size < (*d)[i].size {
						(*d)[i].size -= (*d)[j].size
						(*d)[j].isfile = true
						(*d)[j].id = (*d)[i].id
						j++
						break
					}
				}
				j++
			}
		}
		i--
	}
}

func (d *diskmap) defragWithoutSplit() {
	i := len(*d) - 1
	for i >= 0 {
		if (*d)[i].isfile {
			for j := 0; j < i; j++ {
				if !(*d)[j].isfile {
					if (*d)[j].size == (*d)[i].size {
						(*d)[j].isfile = true
						(*d)[j].id = (*d)[i].id
						(*d)[i].isfile = false
						break
					} else if (*d)[j].size > (*d)[i].size {
						(*d)[j].size -= (*d)[i].size
						(*d).Insert(j, content{id: (*d)[i].id, isfile: true, size: (*d)[i].size})
						(*d)[i+1].isfile = false
						break
					}
				}
			}
		}
		i--
	}
}

func (d *diskmap) Checksum() int {
	count := 0
	running := 0
	for _, c := range *d {
		for i := 0; i < c.size; i++ {
			if c.isfile {
				count += running * c.id
			}
			running++
		}

	}
	return count
}

func loadFile(filename string) []diskmap {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var res []diskmap

	for scanner.Scan() {
		d := diskmap{}
		fileid := 0
		for i, c := range scanner.Text() {
			n, _ := strconv.Atoi(string(c))
			c := content{id: fileid, isfile: i%2 == 0, size: n}
			if c.isfile {
				fileid++
			}
			d = append(d, c)
		}
		res = append(res, d)
	}
	return res
}

func part1(dd []diskmap) int {
	result := 0
	for _, d := range dd {
		fmt.Println(d)
		d.defrag()
		fmt.Println(d)
		checksum := d.Checksum()
		fmt.Println(checksum)
		result += checksum
	}
	return result
}

func part2(dd []diskmap) int {
	result := 0
	for _, d := range dd {
		fmt.Println(d)
		d.defragWithoutSplit()
		fmt.Println(d)
		checksum := d.Checksum()
		fmt.Println(checksum)
		result += checksum
	}
	return result
}

func main() {
	fmt.Println("Part 1 Example:", part1(loadFile("example.txt")))
	//fmt.Println("Part 1:", part1(loadFile("input.txt")))
	/*fmt.Println("Part 2 Example:", part2(loadFile("example.txt")))
	fmt.Println("Part 2:", part2(loadFile("input.txt")))*/
}
