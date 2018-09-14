package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
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

func main() {
	pwd, _ := os.Getwd()
	htm, _ := ioutil.ReadFile(pwd + "/src/templates/second.html")

	doc, err := html.Parse(bytes.NewReader(htm))
	if err != nil {
		fmt.Println("error Parsing HTM", err)
	}

	anchorSlice, err := getLinks(doc)
	if err != nil {
		fmt.Println("Could not get links", err)
	}

	for _, a := range anchorSlice {
		anchor := renderNode(a)
		title := renderNode(a.FirstChild)
		//TODO: find a way to get the Href for the Link struct
		fmt.Println("anchor ", anchor)
		fmt.Println("title ", title)
	}

	fmt.Println("anchorSlice ", anchorSlice)
}
