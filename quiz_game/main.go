package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {

	correct := 0

	filename := flag.String("filename", "problems.csv", "The file with the questions.")
	timeLimit := flag.Int("timelimit", 30, "The time limit to solve the quiz.")
	shuffle := flag.String("shuffle", "n", "To shuffle the question's order.")
	flag.Parse()

	problemsCSV, err := os.Open(*filename)
	check(err)
	defer problemsCSV.Close()

	reader := csv.NewReader(problemsCSV)
	problemsAll, err := reader.ReadAll()
	check(err)

	timerC := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemsLoop:
	for _, p := range makeProblems(problemsAll, *shuffle) {
		answerSet := make(chan string)
		go func() {
			fmt.Printf("%s?: ", p.question)
			userInput := bufio.NewReader(os.Stdin)
			answer, _ := userInput.ReadString('\n')
			answer = strings.ToLower(answer)
			answer = strings.TrimSpace(answer)
			answerSet <- answer
		}()

		select {
		case <-timerC.C:
			fmt.Println("\nTime finished.")
			break problemsLoop
		case answer := <-answerSet:
			if answer == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("Total: %d, Correct: %d.", len(problemsAll), correct)
}

func makeProblems(problems [][]string, shuffle string) []problem {
	result := []problem{}
	if shuffle == string('y') {
		problems = shuffleArray(problems)
	}
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

func shuffleArray(a [][]string) [][]string {
	rand.Seed(time.Now().UnixNano())
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}
