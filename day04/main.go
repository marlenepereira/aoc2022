package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput() ([][]int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	var sections [][]int
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, ",", "-")
		ranges := strings.Split(line, "-")

		var pair []int
		for _, v := range ranges {
			value, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			pair = append(pair, value)
			if len(pair) == 2 {
				sections = append(sections, pair)
				pair = nil
			}
		}
	}

	return sections, nil
}

func overlaps(sections [][]int) (int, int) {
	var partialOverlaps, fullOverlaps int
	for i := 0; i < len(sections); i += 2 {
		// sort if required
		if sections[i][0] > sections[i+1][0] {
			sections[i], sections[i+1] = sections[i+1], sections[i]
		}

		startOne, startTwo := sections[i][0], sections[i+1][0]
		endOne, endTwo := sections[i][1], sections[i+1][1]
		if (startOne >= startTwo && endOne <= endTwo) ||
			(startOne <= startTwo && endOne >= endTwo) {
			fullOverlaps++
			partialOverlaps++
		} else if endOne >= startTwo {
			partialOverlaps++
		}
	}
	return partialOverlaps, fullOverlaps
}

func main() {
	sections, err := readInput()
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	partial, full := overlaps(sections)
	fmt.Printf("part one: %v\n", full)    // 567
	fmt.Printf("part two: %v\n", partial) // 907
}
