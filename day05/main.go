package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const (
	model9000 = iota
	model9001
)

type crane struct {
	model int
	job   *craneJob
}

type craneJob struct {
	procedures []procedure
	supplies   supplies
}

func (c *crane) rearrange() {
	for _, step := range c.job.procedures {
		from := c.job.supplies.stackLocation(step.from)
		to := c.job.supplies.stackLocation(step.to)

		var crates []string
		if c.model == model9000 {
			for i := 0; i < step.moves && len(c.job.supplies[from].crates) != 0; i++ {
				crates = append(crates, c.job.supplies[from].unload()...)
			}
		} else {
			crates = c.job.supplies[from].unload(step.moves)
		}
		c.job.supplies[to].load(crates...)
	}
}

func (c *crane) getFinalArrangement() string {
	var topCrates string
	for _, st := range c.job.supplies {
		topCrates += st.unload()[0]
	}
	return topCrates
}

func (c *crane) assignJob(procedures []procedure, s supplies) {
	var suppliesCopy supplies
	for _, st := range s {
		var crates []string
		crates = append(crates, st.crates...)
		suppliesCopy = append(suppliesCopy, stack{crates})
	}
	c.job = &craneJob{
		procedures: procedures,
		supplies:   suppliesCopy,
	}
}

type procedure struct {
	moves int
	from  int
	to    int
}

type supplies []stack

type stack struct {
	crates []string
}

func (s *supplies) stackLocation(id int) int {
	return id - 1
}

func (s *stack) load(crate ...string) {
	s.crates = append(s.crates, crate...)
}

func (s *stack) unload(q ...int) []string {
	length := len(s.crates)
	if length == 0 {
		return nil
	}

	end := 1
	if len(q) != 0 {
		end = q[0]
	}

	indexRange := length - end
	if indexRange < 0 {
		indexRange = 0
	} else {

	}
	unloaded := s.crates[indexRange:]
	s.crates = s.crates[:indexRange]
	return unloaded
}

func readInputTwo() (supplies, []procedure, error) {
	regex := regexp.MustCompile(`[A-Z-0-9]+`)
	scanner := bufio.NewScanner(os.Stdin)

	// read supplies
	var rawStack stack
	for scanner.Scan() && len(scanner.Bytes()) != 0 {
		line := scanner.Bytes()
		rawStack.load(string(line))
	}

	// parse supplies
	stackIDs := rawStack.unload()
	ids := regex.FindAllString(stackIDs[0], -1)
	stackSupplies := make(supplies, len(ids))

	var stack stack
	windowSize := 4
	for len(rawStack.crates) != 0 {
		row := rawStack.unload()[0]
		idx, start, end := 0, 0, 0
		for end <= len(row) {
			window := row[start:end]
			crate := regex.FindAllString(window, -1)
			if len(crate) != 0 {
				stack = stackSupplies[idx]
				stack.crates = append(stack.crates, crate[0])
				stackSupplies[idx] = stack
			}
			if end-start >= windowSize || len(crate) != 0 {
				start = end
				idx++
			}
			end++
		}
	}

	// read and parse procedures
	var procedures []procedure
	var procedure procedure
	for scanner.Scan() {
		movesString := regex.FindAllString(scanner.Text(), -1)
		var moves []int
		for _, mv := range movesString {
			m, err := strconv.Atoi(mv)
			if err != nil {
				return nil, nil, err
			}
			moves = append(moves, m)
		}
		procedure.moves = moves[0]
		procedure.from = moves[1]
		procedure.to = moves[2]
		procedures = append(procedures, procedure)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return stackSupplies, procedures, nil
}

func main() {
	supplies, procedures, err := readInputTwo()
	if err != nil {
		fmt.Println(err)
	}

	crane9000 := crane{model: model9000}
	crane9000.assignJob(procedures, supplies)
	crane9000.rearrange()

	crane9001 := crane{model: model9001}
	crane9001.assignJob(procedures, supplies)
	crane9001.rearrange()

	fmt.Println(crane9000.getFinalArrangement()) // JCMHLVGMG
	fmt.Println(crane9001.getFinalArrangement()) // LVMRWSSPZ
}
