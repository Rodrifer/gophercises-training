package linkparser

import (
	"os"

	"golang.org/x/net/html"
)

// Link defines a common link structure
type Link struct {
	Href string
	Text string
}

// GetLinks returns a collection of links
func GetLinks(file string) []Link {

	var result []Link
	textLink := ""

	htmlFile, err := os.Open(file)
	check(err)

	doc, err := html.Parse(htmlFile)
	check(err)

	var nt func(*html.Node)
	nt = func(n *html.Node) {
		if n.Type == html.TextNode {
			textLink += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			nt(c)
		}
	}

	var proc func(*html.Node)
	proc = func(n *html.Node) {
		textLink = ""
		nt(n)
		tmpHref := n.Attr[0].Val
		tmpText := textLink
		tmpLink := Link{
			Href: tmpHref,
			Text: tmpText,
		}
		result = append(result, tmpLink)
	}

	var ne func(*html.Node)
	ne = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			proc(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ne(c)
		}
	}
	ne(doc)

	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
