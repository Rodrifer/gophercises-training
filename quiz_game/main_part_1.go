package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func main() {

	correct := 0

	filename := flag.String("filename", "problems.csv", "The file with the questions.")
	flag.Parse()

	problemsCSV, err := os.Open(*filename)
	check(err)

	reader := csv.NewReader(problemsCSV)
	problemsAll, err := reader.ReadAll()
	check(err)

	for _, p := range makeProblems(problemsAll) {
		fmt.Printf("%s?: ", p.question)
		userInput := bufio.NewReader(os.Stdin)
		answer, _ := userInput.ReadString('\n')
		answer = strings.ToLower(answer)
		answer = strings.TrimSpace(answer)
		if p.answer == answer {
			correct++
		}
	}
	fmt.Printf("Total: %d, Correct: %d.", len(problemsAll), correct)
}

func makeProblems(problems [][]string) []problem {
	result := []problem{}
	for _, p := range problems {
		tmp := problem{
			question: p[0],
			answer:   strings.TrimSpace(p[1]),
		}
		result = append(result, tmp)
	}
	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
