package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	correct := 0

	filename := flag.String("filename", "problems.csv", "The file with the questions.")
	flag.Parse()

	problems, err := os.Open(*filename)
	check(err)

	reader := csv.NewReader(problems)
	problemsAll, err := reader.ReadAll()
	check(err)

	for _, p := range problemsAll {
		p[1] = strings.ToLower(p[1])
		p[1] = strings.TrimSpace(p[1])
		fmt.Printf("%s?: ", p[0])
		userInput := bufio.NewReader(os.Stdin)
		answer, _ := userInput.ReadString('\n')
		answer = strings.ToLower(answer)
		answer = strings.TrimSpace(answer)
		if p[1] == answer {
			correct++
		}
	}
	fmt.Printf("Total: %d, Correct: %d.", len(problemsAll), correct)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
