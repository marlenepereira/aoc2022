package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func buildItemPriorities() map[int32]int {
	priority := make(map[int32]int, 51)
	for idx, char := range alphabet {
		priority[char] = idx + 1
	}

	return priority
}

func readInput() ([][]int, error) {
	priorities := buildItemPriorities()
	var rucksacks [][]int

	ch := make([]byte, 1)
	for {
		var rucksack []int
		for {
			if _, err := os.Stdin.Read(ch); err == io.EOF {
				return rucksacks, nil
			}
			if ch[0] == '\n' {
				break
			}

			item := int32(ch[0])
			priority, ok := priorities[item]
			if !ok {
				return nil, fmt.Errorf("couldn't find priority for item: %v", item)
			}
			rucksack = append(rucksack, priority)
		}
		rucksacks = append(rucksacks, rucksack)
	}
}

func getItemPriority(items []int) int {
	middle := len(items) / 2
	sort.Slice(items[:middle], func(i, j int) bool {
		return items[i] < items[j]
	})

	sort.Slice(items[middle:], func(i, j int) bool {
		return items[i+middle] < items[j+middle]
	})

	left, right := 0, middle
	for items[left] != items[right] {
		if items[right] < items[left] {
			right++
		} else if items[right] > items[left] {
			left++
		}
	}
	return items[left]
}

func getGroupBadgePriority(group [][]int) int {
	for _, rucksack := range group {
		sort.Slice(rucksack, func(i, j int) bool {
			return rucksack[i] < rucksack[j]
		})
	}

	elfOne, elfTwo, elfThree := group[0], group[1], group[2]
	var x, y, z int
	for {
		itemOne, itemTwo, itemThree := elfOne[x], elfTwo[y], elfThree[z]
		if itemOne < itemTwo || itemOne < itemThree {
			x++
		} else if itemTwo < itemOne || itemTwo < itemThree {
			y++
		} else if itemThree < itemOne || itemThree < itemTwo {
			z++
		}

		if itemOne == itemTwo && itemOne == itemThree {
			break
		}
	}
	return elfOne[x]
}

func main() {
	rucksacks, err := readInput()
	if err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}

	var group [][]int
	var itemPrioritySum, badgeSum int
	for _, rucksack := range rucksacks {
		itemPrioritySum += getItemPriority(rucksack)
		group = append(group, rucksack)
		if len(group) == 3 {
			badgeSum += getGroupBadgePriority(group)
			group = nil
		}
	}

	fmt.Printf("part one: %v\n", itemPrioritySum) // 7766
	fmt.Printf("part two: %v\n", badgeSum)        // 2415
}
