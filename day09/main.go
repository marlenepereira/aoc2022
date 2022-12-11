package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	R = iota
	L
	U
	D
)

var directionsMap = map[string]int{"R": R, "L": L, "D": D, "U": U}

type rope *list.List

func newRope(knots int) rope {
	r := list.New()
	r.Init()
	for i := 0; i < knots; i++ {
		knot := knot{knotPosition{x: 0, y: 0}}
		r.PushFront(knot)
	}
	return r
}

type knot struct {
	knotPosition
}

type knotPosition struct {
	x int
	y int
}

func readInput() ([]knotPosition, error) {
	scanner := bufio.NewScanner(os.Stdin)

	var positions []knotPosition
	for scanner.Scan() {
		chars := strings.Split(scanner.Text(), " ")
		moves, err := strconv.Atoi(chars[1])
		if err != nil {
			return nil, err
		}

		dir := directionsMap[chars[0]]
		for i := 0; i < moves; i++ {
			pos := knotPosition{x: 0, y: 0}
			switch dir {
			case R:
				pos.x = 1
			case L:
				pos.x = -1
			case U:
				pos.y = 1
			case D:
				pos.y = -1
			}
			positions = append(positions, pos)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return positions, nil
}

func nextMove(val int) int {
	if val > 0 {
		return 1
	} else if val < 0 {
		return -1
	} else {
		return 0
	}
}
func tailNextPosition(head, tail knot) knot {
	yDiff := head.y - tail.y
	xDiff := head.x - tail.x
	if math.Abs(float64(yDiff)) == 2 || math.Abs(float64(xDiff)) == 2 {
		tail.y = tail.y + nextMove(yDiff)
		tail.x = tail.x + nextMove(xDiff)
	}
	return tail
}

func tailVisitedPositions(headNexPositions []knotPosition, rope *list.List) int {
	tailPositions := make(map[knot]int)
	for _, d := range headNexPositions {
		head := rope.Front()
		tail := head.Next()

		headPosition := (head.Value).(knot)
		headPosition.x = headPosition.x + d.x
		headPosition.y = headPosition.y + d.y
		head.Value = headPosition

		var tailPosition knot
		for tail != nil {
			tailPosition = tail.Value.(knot)
			tailPosition = tailNextPosition(headPosition, tailPosition)
			tail.Value = tailPosition
			headPosition = tailPosition
			tail = tail.Next()
		}
		tailPositions[tailPosition]++
	}

	return len(tailPositions)
}

func main() {
	headNextPositions, err := readInput()
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	shortRope := newRope(2)
	longRope := newRope(10)

	fmt.Println(tailVisitedPositions(headNextPositions, shortRope)) // 5695
	fmt.Println(tailVisitedPositions(headNextPositions, longRope))  // 2434
}
