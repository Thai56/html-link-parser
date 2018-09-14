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

func getBody(doc *html.Node, anchorSlice *[]*html.Node, gotAnchors chan bool) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		fmt.Println("DATA ", n.Data)
		if n.Type == html.ElementNode && n.Data == "body" {
			b = n
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			*anchorSlice = append(*anchorSlice, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		fmt.Println("ANCHORS! ", anchorSlice)
		gotAnchors <- true
		return b, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func main() {
	pwd, _ := os.Getwd()
	htm, _ := ioutil.ReadFile(pwd + "/src/templates/second.html")

	doc, err := html.Parse(bytes.NewReader(htm))
	if err != nil {
		fmt.Println("error Parsing HTM", err)
	}

	var anchorSlice []*html.Node

	gotAnchors := make(chan bool)
	go getBody(doc, &anchorSlice, gotAnchors)
	<-gotAnchors

	for _, a := range anchorSlice {
		anchor := renderNode(a)
		title := renderNode(a.FirstChild)
		//TODO: find a way to get the Href for the Link struct
		fmt.Println("anchor ", anchor)
		fmt.Println("title ", title)
	}

	fmt.Println("anchorSlice ", anchorSlice)
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
