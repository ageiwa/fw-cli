package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

func readCmd(word *string, files *[]string, byDir *bool, dir *string) {
	for i, arg := range os.Args {
		if i == 1 {
			*word = arg
		} else if i == 2 && arg == "-dir" {
			*byDir = true
		} else if i >= 2 && !*byDir {
			*files = append(*files, arg)
		} else if i == 3 {
			*dir = arg
		}
	}
}

func findWord(pattern *regexp.Regexp, data []byte) (loc [][]int, err error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	nl := 0

	fLoc := [][]int{}

	for scanner.Scan() {
		nl++

		line := scanner.Text()
		loc := pattern.FindAllIndex([]byte(line), -1)

		if len(loc) == 0 {
			continue
		}

		for _, locEntry := range loc {
			startByte := locEntry[0]
			startRune := []rune(line[:startByte])

			fLoc = append(fLoc, []int{nl, len(startRune)})
		}
	}

	if err := scanner.Err(); err != nil {
		return fLoc, err
	}

	return fLoc, nil
}

func handleFile(pattern *regexp.Regexp, filePath string) (loc [][]int, err error) {
	data, err := os.ReadFile(filePath)
	fLoc := [][]int{}

	if err != nil {
		return fLoc, err
	}

	fLoc, err = findWord(pattern, data)

	if err != nil {
		return fLoc, err
	} else if len(loc) == 0 {
		return fLoc, nil
	}

	return fLoc, nil
}

func toDir(pattern *regexp.Regexp, word string, wg *sync.WaitGroup, dir string) {
	dirEntrys, err := os.ReadDir(dir)

	if err != nil {
		fmt.Printf("Couldn't read directory: %s\n", dir)
		return
	}

	for _, dirEntry := range dirEntrys {
		if dirEntry.IsDir() {
			toDir(pattern, word, wg, dir + dirEntry.Name() + "/")
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			filePath := dir + dirEntry.Name()
			loc, err := handleFile(pattern, filePath)

			if err != nil {
				fmt.Printf("Couldn't process the file: %s\n", filePath)
				return
			}

			for _, locEntry := range loc {
				fmt.Printf("Find '%s' on line %d, column %d in %s\n", word, locEntry[0], locEntry[1], filePath)
			}
		}()
	}
}

func main() {
	word := "май"
	files := []string{}
	dir := "/"
	byDir := false

	readCmd(&word, &files, &byDir, &dir)

	pattern := regexp.MustCompile("(?i)" + word)

	wg := sync.WaitGroup{}

	if byDir {
		toDir(pattern, word, &wg, dir)
		wg.Wait()
		return
	}

	wg.Add(len(files))

	for _, file := range files {
		go func() {
			defer wg.Done()

			loc, err := handleFile(pattern, file)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, locEntry := range loc {
				fmt.Printf("Find '%s' on line %d, column %d in %s\n", word, locEntry[0], locEntry[1], file)
			}
		}()
	}

	wg.Wait()
}
