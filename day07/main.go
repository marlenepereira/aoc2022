package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	diskSpace           = 70000000
	updateRequiredSpace = 30000000

	cmd = iota
	dir
	content

	out  = ".."
	root = "/"
)

type stack struct {
	directories []*directory
}

func (s *stack) push(d ...*directory) {
	s.directories = append(s.directories, d...)
}

func (s *stack) pop(q ...int) []*directory {
	var dirs []*directory
	length := len(s.directories)
	if length == 0 {
		return dirs
	}

	end := 1
	if len(q) > 0 {
		end = q[0]
	}

	indexRange := length - end
	if indexRange < 0 {
		indexRange = 0
	}
	dirs = append(dirs, s.directories[indexRange:]...)
	s.directories = s.directories[:indexRange]
	return dirs
}

func commandOrOutputType(in string) int {
	switch in {
	case "$":
		return cmd
	case "dir":
		return dir
	default:
		return content
	}
}

type file struct {
	name string
	size int
}

type directory struct {
	name           string
	files          []file
	subdirectories map[string]*directory
}

func newDirectory(name string) *directory {
	subDir := make(map[string]*directory)
	return &directory{name: name, subdirectories: subDir}
}

func (d *directory) fileSize() int {
	var size int
	for _, f := range d.files {
		size += f.size
	}
	return size
}

type fileSystem struct {
	root *directory
}

func (f *fileSystem) usedSpace(track func(size int)) int {
	var traverse func(dir *directory) int
	traverse = func(dir *directory) int {
		if dir == nil {
			return 0
		}

		var dirSize int
		for _, d := range dir.subdirectories {
			dirSize += traverse(d)
		}

		total := dir.fileSize() + dirSize
		track(total)
		return total
	}

	return traverse(f.root)
}

func readInput() (*fileSystem, error) {
	f := fileSystem{}
	s := stack{}

	scanner := bufio.NewScanner(os.Stdin)

	var currentDir *directory
	var file file
	for scanner.Scan() {
		line := scanner.Text()

		commandOrOutput := strings.Split(line, " ")
		cmdOrOut := commandOrOutputType(commandOrOutput[0])

		// process commands
		if cmdOrOut == cmd && len(commandOrOutput) > 2 {
			commands := strings.Split(line, " ")
			switch commands[2] {
			case out:
				currentDir = s.pop()[0]
			case root:
				if currentDir == nil {
					currentDir = newDirectory(root)
					s.push(currentDir)
					f.root = currentDir
				} else {
					for len(s.directories) != 0 {
						currentDir = s.pop()[0]
					}
				}
			default:
				// $ cd <dir>
				dirName := commands[2]
				s.push(currentDir)
				currentDir = currentDir.subdirectories[dirName]
			}
		}

		// process directory output
		if cmdOrOut == dir {
			name := commandOrOutput[1]
			dir := newDirectory(name)
			currentDir.subdirectories[dir.name] = dir

		}

		// process file output
		if cmdOrOut == content {
			size, err := strconv.Atoi(commandOrOutput[0])
			if err != nil {
				return nil, err
			}
			file.size = size
			file.name = commandOrOutput[1]
			currentDir.files = append(currentDir.files, file)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &f, nil
}

func processFileSystemSizes(f *fileSystem) (int, int) {
	var dirSizes []int
	var dirSizeSum int
	track := func(size int) {
		if size <= 100000 {
			dirSizeSum += size
		}
		dirSizes = append(dirSizes, size)
	}

	usedSpace := f.usedSpace(track)
	unusedSpace := diskSpace - usedSpace
	spaceToFreeUp := updateRequiredSpace - unusedSpace
	var dirToDelete int
	for _, d := range dirSizes {
		if spaceToFreeUp <= d {
			if dirToDelete == 0 {
				dirToDelete = d
			} else {
				dirToDelete = int(math.Min(float64(dirToDelete), float64(d)))
			}
		}
	}

	return dirSizeSum, dirToDelete
}

func main() {
	fileSystem, err := readInput()
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	dirTotal, toDelete := processFileSystemSizes(fileSystem)
	fmt.Printf("part one: %v\n", dirTotal) // 1307902
	fmt.Printf("part two: %v\n", toDelete) // 7068748
}