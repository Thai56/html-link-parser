package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
	//"strings"
	//"github.com/PuerkitoBio/goquery"
)

type Link struct {
	Href  string
	Title string
}

func getLinks(doc *html.Node) ([]*html.Node, error) {
	var b []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			b = append(b, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func getDoc() *html.Node {
	pwd, _ := os.Getwd()
	//TODO: Test with different edge cases
	htm, _ := ioutil.ReadFile(pwd + "/src/templates/second.html")

	doc, err := html.Parse(bytes.NewReader(htm))
	if err != nil {
		fmt.Println("error Parsing HTM", err)
	}

	return doc
}

type Links struct {
	links []Links
}

func formatLinksStruct(anchorSlice []*html.Node) []Link {
	var links []Link
	for _, a := range anchorSlice {
		title := renderNode(a.FirstChild)
		href := a.Attr[0].Val
		links = append(links, Link{href, title})
	}
	return links
}

func main() {
	doc := getDoc()

	anchorSlice, err := getLinks(doc)
	if err != nil {
		fmt.Println("Could not get links", err)
	}

	var links []Link
	links = formatLinksStruct(anchorSlice)

	fmt.Println(links)
}
