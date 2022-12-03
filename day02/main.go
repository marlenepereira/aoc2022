package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	rock = iota
	paper
	scissors

	loose
	draw
	win
)

type strategyGuide interface {
	play(round []string, config *shapesGameConfig) int
}

type shapesGameConfig struct {
	shapesCode map[string]int
	rules      map[int]int
	scores     map[int]int
}

type shapesGame struct {
	config *shapesGameConfig
}

func (g shapesGame) playStrategy(round []string, strategy strategyGuide) int {
	return strategy.play(round, g.config)
}

type shapeStrategy struct {
	strategy map[string]string
}

func (s shapeStrategy) play(round []string, config *shapesGameConfig) int {
	playerOneShape := config.shapesCode[round[0]]
	playerTwoShape := config.shapesCode[s.strategy[round[1]]]

	var outcome int
	switch playerTwoShape {
	case playerOneShape:
		outcome = draw
	case config.rules[playerOneShape]:
		outcome = loose
	default:
		outcome = win
	}

	return config.scores[outcome] + config.scores[playerTwoShape]
}

type outcomeStrategy struct {
	strategy map[string]int
}

func (o outcomeStrategy) play(round []string, config *shapesGameConfig) int {
	playerOneShape := config.shapesCode[round[0]]
	outcome := o.strategy[round[1]]

	var playShape int
	switch outcome {
	case loose:
		playShape = config.rules[playerOneShape]
	case win:
		for shape, beats := range config.rules {
			if beats == playerOneShape {
				playShape = shape
				break
			}
		}
	case draw:
		playShape = playerOneShape
	}

	return config.scores[outcome] + config.scores[playShape]
}

func readInput() ([][]string, error) {
	var rounds [][]string
	ch := make([]byte, 1)
	for {
		var round []string
		for {
			if _, err := os.Stdin.Read(ch); err == io.EOF {
				return rounds, nil
			}

			if ch[0] == '\n' {
				break
			}

			char := string(ch[0])
			if char != " " {
				round = append(round, char)
			}
		}
		rounds = append(rounds, round)
	}
}

func main() {
	rounds, err := readInput()
	if err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}

	gameConfig := shapesGameConfig{
		shapesCode: map[string]int{"A": rock, "B": paper, "C": scissors},
		rules:      map[int]int{rock: scissors, paper: rock, scissors: paper},
		scores:     map[int]int{loose: 0, draw: 3, win: 6, rock: 1, paper: 2, scissors: 3},
	}

	game := shapesGame{&gameConfig}

	shapeStrategyGuide := shapeStrategy{
		strategy: map[string]string{"X": "A", "Y": "B", "Z": "C"}}
	outcomeStrategyGuide := outcomeStrategy{
		strategy: map[string]int{"X": loose, "Y": draw, "Z": win}}

	var totalScoreOne, totalScoreTwo int
	for _, round := range rounds {
		totalScoreOne += game.playStrategy(round, shapeStrategyGuide)
		totalScoreTwo += game.playStrategy(round, outcomeStrategyGuide)
	}

	fmt.Printf("part one %v\n", totalScoreOne) // 12772
	fmt.Printf("part two %v\n", totalScoreTwo) // 11618
}
