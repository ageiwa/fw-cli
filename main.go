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

func findWord(data []byte) {
	for _, b := range data {
		fmt.Println(b)
	}
}

func main() {
	word := "май"
	files := []string{ "text1.txt" }

	readCmd(&word, &files)

	data, err := os.ReadFile(files[0])

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pattern := regexp.MustCompile(word)
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

		fmt.Printf(`Find "%s" on line %d, column %d`, word, nl, len(startRune))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
	}
}