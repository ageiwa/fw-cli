package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func readCmd(word *string, files *[]string) {
	for i, arg := range os.Args {
		if i == 1 {
			*word = arg
		} else if i > 1 {
			*files = append(*files, arg)
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

			fLoc = append(fLoc, []int{ nl, len(startRune) })
		}
	}

	if err := scanner.Err(); err != nil {
		return fLoc, err
	}

	return fLoc, nil
}

func main() {
	word := "май"
	files := []string{ "text1.txt", "text2.txt", "text3.txt" }

	readCmd(&word, &files)

	pattern := regexp.MustCompile("(?i)" + word)

	for _, file := range files {
		data, err := os.ReadFile(file)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		loc, err := findWord(pattern, data)

		if err != nil {
			fmt.Println(err.Error())
			continue
		} else if len(loc) == 0 {
			continue
		}

		for _, locEntry := range loc {
			fmt.Printf("Find '%s' on line %d, column %d in %s\n", word, locEntry[0], locEntry[1], file)
		}
	}
}