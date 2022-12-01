package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
)

type calorieCount []int

func (c calorieCount) Len() int {
	return len(c)
}

func (c calorieCount) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c calorieCount) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *calorieCount) Push(x any) {
	*c = append(*c, x.(int))
}

func (c *calorieCount) Pop() any {
	copyCal := *c
	n := c.Len() - 1
	pop := copyCal[n]
	*c = copyCal[:n]
	return pop
}

func readInput() ([]int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	var calories []int
	var current int
	for scanner.Scan() {
		calString := scanner.Text()

		if calString == "" {
			// found empty line
			calories = append(calories, current)
			current = 0
			continue
		}
		cal, err := strconv.Atoi(calString)
		if err != nil {
			return nil, err
		}
		current += cal
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return calories, nil
}

func main() {
	cal, err := readInput()
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	minHeap := calorieCount{}
	kth := 3
	var maxSum int
	for _, c := range cal {
		maxSum += c
		heap.Push(&minHeap, c)
		for minHeap.Len() > kth {
			pop := heap.Pop(&minHeap)
			maxSum -= pop.(int)
		}
	}

	var maxCount int
	for minHeap.Len() != 0 {
		maxCount = heap.Pop(&minHeap).(int)
	}

	fmt.Printf("part one: %v\n", maxCount)
	fmt.Printf("part two: %v\n", maxSum)
}
