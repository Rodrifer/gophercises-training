package main

import (
	"flag"
	"fmt"
	"gophercises-training/link_parser/linkparser"
)

func main() {

	file := flag.String("html", "ex1.html", "The HTML file to get the Links")
	flag.Parse()

	links := linkparser.GetLinks(*file)
	fmt.Println(links)
}
