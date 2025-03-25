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

func findWord(pattern *regexp.Regexp, data []byte) (loc []int, err error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	nl := 0

	for scanner.Scan() {
		nl++

		line := scanner.Text()
		loc := pattern.FindIndex([]byte(line))

		if len(loc) == 0 {
			continue
		}

		startByte := loc[0]
		startRune := []rune(line[:startByte])

		return []int{ nl, len(startRune) }, nil
	}

	if err := scanner.Err(); err != nil {
		return loc, err
	}

	return loc, nil
}

func main() {
	word := "май"
	files := []string{ "text1.txt", "text2.txt", "text3.txt" }

	readCmd(&word, &files)

	pattern := regexp.MustCompile(word)

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

		fmt.Printf("Find '%s' on line %d, column %d in %s\n", word, loc[0], loc[1], file)
	}
}