package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	north = iota
	south
	east
	west
)

func readInput() (forest, error) {
	var grid forest
	for {
		ch := make([]byte, 1)
		var row []int
		for {
			if _, err := os.Stdin.Read(ch); err == io.EOF {
				return grid, nil
			}

			if ch[0] == '\n' {
				grid = append(grid, row)
				break
			}

			val, err := strconv.Atoi(string(ch[0]))
			if err != nil {
				return nil, err
			}
			row = append(row, val)
		}
	}
}

type forest [][]int

func (g forest) walk(direction int, fromTree int, boundary int, fixed int, start int) (int, bool) {
	var treesCount int
	isVisible := true
	var neighbourTree int
	for start != boundary {
		switch direction {
		case north:
			start--
			neighbourTree = g[start][fixed]
		case south:
			start++
			neighbourTree = g[start][fixed]
		case east:
			start++
			neighbourTree = g[fixed][start]
		case west:
			start--
			neighbourTree = g[fixed][start]
		}

		treesCount++
		if neighbourTree >= fromTree {
			isVisible = false
			return treesCount, isVisible
		}
	}

	return treesCount, isVisible
}

func visibleTreeCount(g forest) (int, int) {
	northEdge, southEdge := 0, len(g)-1
	westEdge, eastEdge := 0, len(g[0])-1
	var count, treeCountMax int
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			fromTree := g[y][x]
			northCount, northVisible := g.walk(north, fromTree, northEdge, x, y)
			southCount, southVisible := g.walk(south, fromTree, southEdge, x, y)
			eastCount, eastVisible := g.walk(east, fromTree, eastEdge, y, x)
			westCount, westVisible := g.walk(west, fromTree, westEdge, y, x)

			if northVisible || southVisible || eastVisible || westVisible {
				count++
			}

			currentCount := northCount * southCount * eastCount * westCount
			treeCountMax = int(math.Max(float64(currentCount), float64(treeCountMax)))
		}
	}

	return count, treeCountMax
}

func main() {
	grid, err := readInput()
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	visible, trees := visibleTreeCount(grid)

	fmt.Printf("part one: %v\n", visible) // 1803
	fmt.Printf("part two: %v\n", trees)   // 268912
}
