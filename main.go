package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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

func handleFile(pattern *regexp.Regexp, word string, filePath string) (res []string, err error) {
	data, err := os.ReadFile(filePath)
	resList := []string{}

	if err != nil {
		return resList, err
	}

	loc, err := findWord(pattern, data)

	if err != nil {
		return resList, err
	} else if len(loc) == 0 {
		return resList, nil
	}

	for _, locEntry := range loc {
		result := fmt.Sprintf("Find '%s' on line %d, column %d in %s\n", word, locEntry[0], locEntry[1], filePath)
		resList = append(resList, result)
	}

	return resList, nil
}

func main() {
	word := "май"
	files := []string{"text1.txt", "text2.txt", "text3.txt"}
	dir := "/"
	byDir := false

	readCmd(&word, &files, &byDir, &dir)

	pattern := regexp.MustCompile("(?i)" + word)

	for _, file := range files {
		handleFile(pattern, word, file)
	}
}
