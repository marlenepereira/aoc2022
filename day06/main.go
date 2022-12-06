package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func receiveStream(markSize int) (int, error) {
	input, err := os.Open("day06/input/input.txt")
	if err != nil {
		return 0, err
	}

	defer input.Close()

	signals := make(map[string]int)
	var stream []byte
	var start, end, duplicates int
	for {
		ch := make([]byte, 1)
		if _, err := input.Read(ch); err == io.EOF {
			return 0, nil
		}

		stream = append(stream, ch...)
		end++
		signals[string(ch[0])]++
		count, _ := signals[string(ch[0])]
		if count > 1 {
			duplicates++
		}

		for len(stream[start:end]) > markSize {
			toRemove := string(stream[start])
			signals[toRemove]--
			count, _ := signals[toRemove]
			if count >= 1 {
				duplicates--
			}
			start++
		}

		if len(stream[start:end]) == markSize && duplicates == 0 {
			break
		}
	}

	return len(stream), nil
}

func main() {
	signalOne, err := receiveStream(4)
	if err != nil {
		log.Fatalf("Error receiving stream: %v", err)
	}

	signalTwo, err := receiveStream(14)
	if err != nil {
		log.Fatalf("Error receiving stream: %v", err)
	}

	fmt.Printf("part one: %v\n", signalOne) // 1134
	fmt.Printf("part two: %v\n", signalTwo) // 2263
}
