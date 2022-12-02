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
	// create a max heap
	return c[i] > c[j]
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

func readInput() (*calorieCount, error) {
	scanner := bufio.NewScanner(os.Stdin)

	calories := &calorieCount{}
	var current int
	for scanner.Scan() {
		calString := scanner.Text()

		if calString == "" {
			// found empty line
			heap.Push(calories, current)
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

	// push last calorie count
	heap.Push(calories, current)
	return calories, nil
}

func main() {
	cal, err := readInput()
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	maxCount := heap.Pop(cal).(int)
	maxSum := maxCount
	for k := 2; k > 0; k-- {
		pop := heap.Pop(cal).(int)
		fmt.Println(pop)
		maxSum += pop
	}

	fmt.Printf("part one: %v\n", maxCount) // 69289
	fmt.Printf("part two: %v\n", maxSum)   // 205615
}
